package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/tigrisdata-community/glue/internal/store"
	"github.com/tigrisdata-community/glue/web/discourse"
)

var (
	//go:embed massage-system-prompt.txt
	cleanupSystemPrompt string
)

type DiscourseQuestion struct {
	Title string          `json:"title"`
	Slug  string          `json:"slug"`
	Posts []DiscoursePost `json:"posts"`
}

type DiscoursePost struct {
	Body     string `json:"body"`
	UserID   string `json:"userID"`
	Accepted bool   `json:"accepted"`
}

func discourseMassage(ctx context.Context) error {
	st, err := store.NewS3API(ctx, *storeBucket)
	if err != nil {
		return err
	}

	discourseTopics := store.JSON[discourse.TopicResult]{
		Underlying: st,
		Prefix:     "discourse",
	}

	discourseThreads := store.JSON[DiscourseQuestion]{
		Underlying: st,
		Prefix:     "discourse-thread",
	}

	ai := openai.NewClient(
		option.WithAPIKey(*openAIAPIKey),
		option.WithBaseURL(*openAIAPIBase),
	)
	_ = ai

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

		for i, post := range topic.PostStream.Posts {
			if post.Username == "system" {
				continue
			}

			params := openai.ChatCompletionNewParams{
				Messages: []openai.ChatCompletionMessageParamUnion{
					openai.SystemMessage(cleanupSystemPrompt),
					openai.UserMessage(post.Cooked),
				},
			}

			resp, err := ai.Chat.Completions.New(ctx, params)
			if err != nil {
				errs = append(errs, fmt.Errorf("while censoring the %d message in %s: %w", i, k, err))
				continue
			}

			thread.Posts = append(thread.Posts, DiscoursePost{
				Body:     resp.Choices[0].Message.Content,
				UserID:   fmt.Sprint(post.UserTitle, " ", post.UserID),
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
