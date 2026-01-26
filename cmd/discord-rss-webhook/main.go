package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/facebookgo/flagenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pstuifzand/ekster/pkg/jsonfeed"
	"github.com/tigrisdata-community/glue/internal"
	"github.com/tigrisdata-community/glue/internal/store"
	"github.com/tigrisdata-community/glue/web"
	"github.com/tigrisdata-community/glue/web/discordwebhook"
	"github.com/tigrisdata-community/glue/web/useragent"
)

var (
	discordAvatarURL  = flag.String("discord-avatar-url", "https://gtm-glue-discord-webhook.t3.storage.dev/avatars/ty.webp", "Discord pseudo-user avatar URL")
	discordUsername   = flag.String("discord-username", "Ty", "Discord pseudo-user username")
	discordWebhookURL = flag.String("discord-webhook-url", "", "Discord webhook URL")
	feedURL           = flag.String("feed-url", "https://www.tigrisdata.com/blog/feed.json", "Blog JSONfeed")
	storeBucket       = flag.String("store-bucket", "", "The Tigris bucket used to store data")
)

type SeenURL struct {
}

func main() {
	flagenv.Parse()
	flag.Parse()

	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(errors.New("main exited"))

	slog.Info(
		"starting up",
		"has-discord-avatar-url", *discordAvatarURL != "",
		"discord-username", *discordUsername,
		"has-discord-webhook-url", *discordWebhookURL != "",
		"store-bucket", *storeBucket,
		"args", flag.Args(),
	)

	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "fatal:", err)
	}
}

func run(ctx context.Context) error {
	st, err := store.NewS3API(ctx, *storeBucket)
	if err != nil {
		return err
	}

	seenURLs := &store.JSON[string]{
		Underlying: st,
		Prefix:     "seen-urls",
	}

	ua := useragent.Generate("tigris-gtm-glue", "https://tigrisdata.com")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, *feedURL, nil)
	if err != nil {
		return fmt.Errorf("can't make request: %w", err)
	}
	req.Header.Set("User-Agent", ua)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("can't fetch response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return web.NewError(http.StatusOK, resp)
	}

	defer resp.Body.Close()
	feed, err := jsonfeed.Parse(resp.Body)
	if err != nil {
		return fmt.Errorf("can't parse jsonfeed: %w", err)
	}

	slog.Info("got feed", "title", feed.Title)

	var errs []error

	for _, item := range feed.Items {
		key := internal.SHA256sum(item.ID)
		err := seenURLs.Exists(ctx, key)

		if err != nil && !errors.Is(err, store.ErrNotFound) {
			slog.Error("can't verify item exists in store", "err", err)
			errs = append(errs, err)
			continue
		}

		if err != nil && errors.Is(err, store.ErrNotFound) {
			// do Discord egress

			req := discordwebhook.Send(*discordWebhookURL, discordwebhook.Webhook{
				Username:  *discordUsername,
				AvatarURL: *discordAvatarURL,
				Content:   fmt.Sprintf("New blogpost: %s", item.URL),
			})
			req.Header.Set("User-Agent", ua)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				slog.Error("can't egress discord webhook", "err", err)
				errs = append(errs, err)
				continue
			}

			if err := discordwebhook.Validate(resp); err != nil {
				slog.Error("can't validate discord webhook response", "err", err)
				errs = append(errs, fmt.Errorf("can't post webhook: %w", err))
				continue
			}
		}

		slog.Info("seen item", "key", key, "title", item.Title, "id", item.ID, "summary", item.Summary)
		if err := seenURLs.Set(ctx, key, item.Title); err != nil {
			slog.Error("can't store item info in store", "err", err)
			errs = append(errs, err)
			continue
		}
	}

	if len(errs) != 0 {
		return fmt.Errorf("got errors processing feed items:\n%w", errors.Join(errs...))
	}

	return nil
}
