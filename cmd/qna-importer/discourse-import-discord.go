package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-faker/faker/v4"
	"github.com/tigrisdata-community/glue/internal/store"
	"github.com/tigrisdata-community/glue/web/discordwebhook"
	"github.com/tigrisdata-community/glue/web/sdcpp"
	"github.com/tigrisdata-community/glue/web/useragent"
	"github.com/tigrisdata/storage-go"
)

var (
	discordToken      = flag.String("discord-token", "", "Discord bot token")
	discordWebhookURL = flag.String("discord-webhook-url", "", "Discord webhook URL")
	sdcppURL          = flag.String("sdcpp-url", "", "stable-diffusion.cpp server URL")

	postDelay = flag.Duration("post-delay", 5*time.Second, "delay between post creation attempts")
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

	// discourse key -> discord channel ID
	discourseToDiscord := store.JSON[string]{
		Underlying: st,
		Prefix:     "discord-thread-mapping",
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

	// // For testing, comment out in prod
	// threads = append([]string{}, threads[0])

	delayTick := time.NewTicker(*postDelay)
	defer delayTick.Stop()

	for _, key := range threads {
		lg := slog.With("key", key)
		thread, err := discourseThreads.Get(ctx, key)
		if err != nil {
			lg.Error("can't fetch thread", "err", err)
			errs = append(errs, fmt.Errorf("while fetching thread %s: %w", key, err))
			continue
		}

		if len(thread.Posts) == 0 {
			lg.Debug("skipping thread with no posts")
			slog.Info("skipping thread with no posts", "slug", thread.Slug, "title", thread.Title)
			continue
		}

		if discordID, err := discourseToDiscord.Get(ctx, key); err == nil {
			lg.Info("skipping thread we've already saved", "thread", key, "discord_id", discordID, "err", err)
			continue
		}

		op := thread.Posts[0]
		user := ug.Get(ctx, op.UserID)
		wh := discordwebhook.Webhook{
			Content:    op.Body,
			ThreadName: thread.Title,
			Username:   user.Username,
		}

		if user.AvatarKey != "" {
			wh.AvatarURL = fmt.Sprintf("https://%s.t3.storage.dev/%s", *storeBucket, user.AvatarKey)
		}

		q := u.Query()
		q.Del("thread_id")
		q.Set("wait", "true")

		u.RawQuery = q.Encode()

		<-delayTick.C
		req := discordwebhook.Send(u.String(), wh)
		req.Header.Set("User-Agent", useragent.Generate("tigris-gtm-glue", "https://tigrisdata.com"))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			lg.Error("can't create thread", "err", err)
			errs = append(errs, fmt.Errorf("can't create thread %s: %w", key, err))
			continue
		}
		tcr, err := discordwebhook.ParseThreadCreation(resp)
		if err != nil {
			lg.Error("can't parse thread creation response", "err", err)
			errs = append(errs, fmt.Errorf("can't parse thread creation response for %s: %w", key, err))
			continue
		}

		discourseToDiscord.Set(ctx, key, tcr.ChannelID)

		q = u.Query()
		q.Del("wait")
		q.Set("thread_id", tcr.ChannelID)

		u.RawQuery = q.Encode()

		whurl := u.String()

		slog.Info("created discord forum thread", "slug", thread.Slug, "title", thread.Title, "id", tcr.ChannelID)

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

			<-delayTick.C
			req := discordwebhook.Send(whurl, wh)
			req.Header.Set("User-Agent", useragent.Generate("tigris-gtm-glue", "https://tigrisdata.com"))
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				errs = append(errs, fmt.Errorf("can't send %s %dth reply: %w", key, i, err))
			}

			if err := discordwebhook.Validate(resp); err != nil {
				errs = append(errs, fmt.Errorf("can't post webhook: %w", err))
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
