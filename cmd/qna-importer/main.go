package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"log/slog"

	"github.com/facebookgo/flagenv"
	_ "github.com/joho/godotenv/autoload"
)

var (
	discourseURL    = flag.String("discourse-url", "https://community.fly.io", "Base Discourse URL")
	discourseTagURL = flag.String("discourse-tag-url", "/tags/c/questions-and-help/11/tigris.json", "Discourse server URL")
	openAIAPIBase   = flag.String("openai-api-base", "", "OpenAI API base URL")
	openAIAPIKey    = flag.String("openai-api-key", "", "OpenAI API key")
	openAIModel     = flag.String("openai-model", "gpt-oss-120b", "OpenAI model")
	storeBucket     = flag.String("store-bucket", "", "The Tigris bucket used to store data")
)

func main() {
	flagenv.Parse()
	flag.Parse()

	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(errors.New("main exited"))

	slog.Info(
		"starting up",
		"has-discord-token", *discordToken != "",
		"has-discord-webhook-url", *discordWebhookURL != "",
		"discourse-url", *discourseURL,
		"discourse-tag-url", *discourseTagURL,
		"openai-api-base", *openAIAPIBase,
		"has-openai-api-key", *openAIAPIKey != "",
		"openai-model", *openAIModel,
		"store-bucket", *storeBucket,
		"args", flag.Args(),
	)

	switch flag.Arg(0) {
	case "discourse-scrape":
		if err := discourseScrape(ctx); err != nil {
			log.Fatal("error:", err)
		}

	case "discourse-massage":
		if err := discourseMassage(ctx); err != nil {
			log.Fatal("error:", err)
		}

	case "discord-import-discord":
		if err := discourseImportDiscord(ctx); err != nil {
			log.Fatal("error:", err)
		}

	default:
		log.Fatalf("ERROR unknown command: %q", flag.Arg(0))
	}
}
