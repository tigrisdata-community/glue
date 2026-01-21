package main

import (
	"context"
	"fmt"

	"github.com/tigrisdata-community/glue/internal/store"
	"github.com/tigrisdata-community/glue/web/discourse"
)

func discourseScrape(ctx context.Context) error {
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
