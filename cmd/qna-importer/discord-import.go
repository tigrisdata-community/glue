package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/tigrisdata-community/glue/internal/store"
)

var (
	discordForumChannel = flag.String("discord-forum-channel", "1457749835871686737", "Discord forum channel to operate in")
	discordGuild        = flag.String("discord-guild", "1457741299041046581", "Discord guild to operate in")
	discordToken        = flag.String("discord-token", "", "Discord bot token")
)

func discordImport(ctx context.Context) error {
	st, err := store.NewS3API(ctx, *storeBucket)
	if err != nil {
		return err
	}

	discourseThreads := store.JSON[DiscourseQuestion]{
		Underlying: st,
		Prefix:     "discourse-thread",
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
		_, err = dc.ForumThreadStartComplex(*discordForumChannel,
			&discordgo.ThreadStart{
				Name: thread.Title,
			},
			msgSend,
		)
		if err != nil {
			errs = append(errs, fmt.Errorf("while creating discord thread for %s: %w", key, err))
			continue
		}

		slog.Info("created discord forum thread", "slug", thread.Slug, "title", thread.Title)
	}

	if len(errs) != 0 {
		return fmt.Errorf("got errors during import: %w", errors.Join(errs...))
	}

	return nil
}
