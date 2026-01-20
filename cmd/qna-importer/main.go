package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"

	"github.com/facebookgo/flagenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/tigrisdata-community/glue/internal/store"
	"github.com/tigrisdata-community/glue/web/discourse"
)

var (
	discourseURL    = flag.String("discourse-url", "https://community.fly.io", "Base Discourse URL")
	discourseTagURL = flag.String("discourse-tag-url", "/tags/c/questions-and-help/11/tigris.json", "Discourse server URL")
	storeBucket     = flag.String("store-bucket", "", "The Tigris bucket used to store data")
)

type DiscourseQuestion struct {
	Title string          `json:"title"`
	Slug  string          `json:"slug"`
	Posts []DiscoursePost `json:"posts"`
}

type DiscoursePost struct {
	Body     string `json:"body"`
	Accepted bool   `json:"accepted"`
}

func main() {
	flagenv.Parse()
	flag.Parse()

	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(errors.New("main exited"))

	slog.Info(
		"starting up",
		"discourse-url", *discourseURL,
		"discourse-tag-url", *discourseTagURL,
		"store-bucket", *storeBucket,
		"args", flag.Args(),
	)

	switch flag.Arg(0) {
	case "discourse-import":
		if err := discourseImport(ctx); err != nil {
			log.Fatal("error:", err)
		}

	case "discourse-massage":
		if err := discourseMassage(ctx); err != nil {
			log.Fatal("error:", err)
		}

	default:
		log.Fatalf("ERROR unknown command: %q", flag.Arg(0))
	}
}

func discourseMassage(ctx context.Context) error {
	st, err := store.NewS3API(ctx, *storeBucket)
	if err != nil {
		return err
	}
	_ = st

	discourseTopics := store.JSON[discourse.TopicResult]{
		Underlying: st,
		Prefix:     "discourse",
	}
	_ = discourseTopics

	discourseThreads := store.JSON[DiscourseQuestion]{
		Underlying: st,
		Prefix:     "discourse-thread",
	}
	_ = discourseThreads

	keys, err := discourseTopics.List(ctx, "")
	if err != nil {
		return fmt.Errorf("can't list cached topics: %w", err)
	}

	var errs []error

	for _, k := range keys {
		fmt.Println(k)

		topic, err := discourseTopics.Get(ctx, k)
		if err != nil {
			errs = append(errs, fmt.Errorf("while fetching %s: %w", k, err))
			continue
		}

		thread := DiscourseQuestion{
			Title: topic.Title,
			Slug:  k,
		}

		for _, post := range topic.PostStream.Posts {
			if post.Username == "system" {
				continue
			}

			thread.Posts = append(thread.Posts, DiscoursePost{
				Body:     post.Cooked,
				Accepted: post.AcceptedAnswer,
			})
		}

		if err := discourseThreads.Set(ctx, k, thread); err != nil {
			errs = append(errs, fmt.Errorf("while setting thread for %s: %w", k, err))
			continue
		}
	}

	if len(errs) != 0 {
		return fmt.Errorf("got a few errors: %w", errors.Join(errs...))
	}

	return nil
}

func discourseImport(ctx context.Context) error {
	st, err := store.NewS3API(ctx, *storeBucket)
	if err != nil {
		return err
	}
	_ = st

	discourseTopics := store.JSON[discourse.TopicResult]{
		Underlying: st,
		Prefix:     "discourse",
	}

	catr, err := discourse.GetCategoryAndTag(ctx, *discourseURL+*discourseTagURL)
	if err != nil {
		return err
	}

	for _, topic := range catr.TopicList.Topics {
		fmt.Println(topic.JSONURL())

		topicData, err := discourse.GetTopic(ctx, *discourseURL+topic.JSONURL())
		if err != nil {
			return err
		}

		if err := discourseTopics.Set(ctx, fmt.Sprintf("%d-%s", topic.ID, topic.Slug), *topicData); err != nil {
			return err
		}
	}

	return nil
}
