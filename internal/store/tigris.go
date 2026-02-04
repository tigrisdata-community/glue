package store

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/tigrisdata/storage-go/simplestorage"
)

func NewTigris(ctx context.Context, bucket string) (Interface, error) {
	client, err := simplestorage.New(ctx,
		simplestorage.WithBucket(bucket),
	)

	if err != nil {
		return nil, err
	}

	return &Tigris{cli: client, bucket: bucket}, nil
}

type Tigris struct {
	cli    *simplestorage.Client
	bucket string
}

func (t *Tigris) Delete(ctx context.Context, key string) error {
	if _, err := t.cli.Head(ctx, key); err != nil {
		return fmt.Errorf("%w: %w", ErrNotFound, err)
	}
	iopsMetrics.WithLabelValues("tigris", "HeadObject")

	if err := t.cli.Delete(ctx, key); err != nil {
		return err
	}
	iopsMetrics.WithLabelValues("tigris", "DeleteObject")

	return nil
}

func (t *Tigris) Exists(ctx context.Context, key string) error {
	if _, err := t.cli.Head(ctx, key); err != nil {
		return fmt.Errorf("%w: %w", ErrNotFound, err)
	}
	iopsMetrics.WithLabelValues("tigris", "HeadObject")

	return nil
}

func (t *Tigris) Get(ctx context.Context, key string) ([]byte, error) {
	out, err := t.cli.Get(ctx, key)
	iopsMetrics.WithLabelValues("s3api", "GetObject")
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
	}
	defer out.Body.Close()

	b, err := io.ReadAll(out.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read s3 object: %w", err)
	}
	return b, nil
}

func (t *Tigris) Set(ctx context.Context, key string, value []byte) error {
	if _, err := t.cli.Put(ctx, &simplestorage.Object{
		Key:  key,
		Body: io.NopCloser(bytes.NewBuffer(value)),
		Size: int64(len(value)),
	}); err != nil {
		return err
	}

	iopsMetrics.WithLabelValues("tigris", "PutObject")
	return nil
}

func (t *Tigris) List(ctx context.Context, prefix string) ([]string, error) {
	items, err := t.cli.List(ctx, simplestorage.WithPrefix(prefix))
	if err != nil {
		return nil, fmt.Errorf("can't list items: %w", err)
	}
	iopsMetrics.WithLabelValues("s3api", "ListObjectsV2")

	var result []string

	for _, item := range items.Items {
		result = append(result, item.Key)
	}

	return result, nil
}
