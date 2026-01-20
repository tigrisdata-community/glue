package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/bwmarrin/discordgo"
	"github.com/go-faker/faker/v4"
	"github.com/tigrisdata-community/glue/internal/store"
	"github.com/tigrisdata-community/glue/web"
	"github.com/tigrisdata-community/glue/web/discordwebhook"
	"github.com/tigrisdata-community/glue/web/sdcpp"
	"github.com/tigrisdata-community/glue/web/useragent"
	"github.com/tigrisdata/storage-go"
)

var (
	discordForumChannel = flag.String("discord-forum-channel", "1457749835871686737", "Discord forum channel to operate in")
	discordToken        = flag.String("discord-token", "", "Discord bot token")
	discordWebhookURL   = flag.String("discord-webhook-url", "", "Discord webhook URL")
	sdcppURL            = flag.String("sdcpp-url", "", "stable-diffusion.cpp server URL")
)

func discourseImportDiscord(ctx context.Context) error {
	st, err := store.NewS3API(ctx, *storeBucket)
	if err != nil {
		return err
	}

	discourseThreads := store.JSON[DiscourseQuestion]{
		Underlying: st,
		Prefix:     "discourse-thread",
	}

	tigris, err := storage.New(ctx)
	if err != nil {
		return err
	}

	ug := &UserGenerator{
		Storage: store.JSON[FakeUser]{
			Underlying: st,
			Prefix:     "discord-generated-usernames",
		},
		AvatarGen: &AvatarGen{
			sd: &sdcpp.Client{
				HTTP:      http.DefaultClient,
				APIServer: *sdcppURL,
			},
			tigris: tigris,
			bucket: *storeBucket,
			prefix: "avatars",
		},
	}

	dc, err := discordgo.New("Bot " + *discordToken)
	if err != nil {
		return fmt.Errorf("can't create discord bot client: %w", err)
	}

	if err := dc.Open(); err != nil {
		return fmt.Errorf("can't open discord connection: %w", err)
	}

	defer dc.Close()

	threads, err := discourseThreads.List(ctx, "")
	if err != nil {
		return fmt.Errorf("can't list discourse threads: %w", err)
	}

	u, err := url.Parse(*discordWebhookURL)
	if err != nil {
		return fmt.Errorf("discord webhook URL doesn't parse: %w", err)
	}

	var errs []error

	for _, key := range threads {
		thread, err := discourseThreads.Get(ctx, key)
		if err != nil {
			errs = append(errs, fmt.Errorf("while fetching thread %s: %w", key, err))
			continue
		}

		if len(thread.Posts) == 0 {
			slog.Info("skipping thread with no posts", "slug", thread.Slug, "title", thread.Title)
			continue
		}

		// Create forum thread with the title and first post content
		msgSend := &discordgo.MessageSend{
			Content: thread.Posts[0].Body,
		}

		// Create the thread in the forum channel
		ch, err := dc.ForumThreadStartComplex(*discordForumChannel,
			&discordgo.ThreadStart{
				Name: thread.Title,
			},
			msgSend,
		)
		if err != nil {
			errs = append(errs, fmt.Errorf("while creating discord thread for %s: %w", key, err))
			continue
		}

		q := u.Query()
		q.Set("thread_id", ch.ID)

		u.RawQuery = q.Encode()

		whurl := u.String()

		slog.Info("created discord forum thread", "slug", thread.Slug, "title", thread.Title, "id", ch.ID)

		for i, post := range thread.Posts {
			if i == 0 {
				continue
			}

			user := ug.Get(ctx, post.UserID)
			wh := discordwebhook.Webhook{
				Content:  post.Body,
				Username: user.Username,
			}
			slog.Info("got user", "user", user)

			if user.AvatarKey != "" {
				wh.AvatarURL = fmt.Sprintf("https://%s.t3.storage.dev/%s", *storeBucket, user.AvatarKey)
			}

			req := discordwebhook.Send(whurl, wh)
			req.Header.Set("User-Agent", useragent.Generate("tigris-gtm-glue", "https://tigrisdata.com"))
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				errs = append(errs, fmt.Errorf("can't send %s %dth reply: %w", key, i, err))
			}
			if resp.StatusCode != http.StatusNoContent {
				errs = append(errs, web.NewError(http.StatusNoContent, resp))
			}
		}
	}

	if len(errs) != 0 {
		return fmt.Errorf("got errors during import: %w", errors.Join(errs...))
	}

	return nil
}

type UserGenerator struct {
	Storage   store.JSON[FakeUser]
	AvatarGen *AvatarGen
}

func (ug *UserGenerator) Get(ctx context.Context, key string) FakeUser {
	result, err := ug.Storage.Get(ctx, key)
	if err != nil {
		slog.Debug("got error fetching username", "key", key, "err", err)

		result = FakeUser{
			ActualUID: key,
			Username:  faker.Name(),
		}

		avatarKey, err := ug.AvatarGen.GenerateAndUpload(ctx, key)
		if err != nil {
			slog.Error("can't render and upload avatar", "err", err)
		} else {
			result.AvatarKey = avatarKey
		}

		ug.Storage.Set(ctx, key, result)
	}

	return result
}

type FakeUser struct {
	ActualUID string `json:"actual_uid"`
	Username  string `json:"username"`
	AvatarKey string `json:"avatar_url"`
}
