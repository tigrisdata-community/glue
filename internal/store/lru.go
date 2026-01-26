package store

import (
	"context"
	"fmt"

	lru "github.com/hashicorp/golang-lru/v2"
)

type LRU struct {
	cache      *lru.Cache[string, []byte]
	underlying Interface
}

func NewLRUCache(underlying Interface) (*LRU, error) {
	cache, err := lru.New[string, []byte](512)
	if err != nil {
		return nil, fmt.Errorf("can't create LRU cache: %w", err)
	}

	return &LRU{
		cache:      cache,
		underlying: underlying,
	}, nil
}

func (l *LRU) Delete(ctx context.Context, key string) error {
	l.cache.Remove(key)
	return l.underlying.Delete(ctx, key)
}

func (l *LRU) Exists(ctx context.Context, key string) error {
	return l.underlying.Exists(ctx, key)
}

func (l *LRU) Get(ctx context.Context, key string) ([]byte, error) {
	result, ok := l.cache.Get(key)
	if ok {
		iopsMetrics.WithLabelValues("lru", "cache_read")
		return result, nil
	}

	iopsMetrics.WithLabelValues("lru", "cache_load")
	result, err := l.underlying.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	l.cache.Add(key, result)

	return result, nil
}

func (l *LRU) Set(ctx context.Context, key string, value []byte) error {
	l.cache.Add(key, value)
	return l.underlying.Set(ctx, key, value)
}

func (l *LRU) List(ctx context.Context, prefix string) ([]string, error) {
	return l.underlying.List(ctx, prefix)
}
