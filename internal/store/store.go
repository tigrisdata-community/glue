package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// ErrNotFound is returned when the store implementation cannot find the value
	// for a given key.
	ErrNotFound = errors.New("store: key not found")

	// ErrCantDecode is returned when a store adaptor cannot decode the store format
	// to a value used by the code.
	ErrCantDecode = errors.New("store: can't decode value")

	// ErrCantEncode is returned when a store adaptor cannot encode the value into
	// the format that the store uses.
	ErrCantEncode = errors.New("store: can't encode value")

	// ErrBadConfig is returned when a store adaptor's configuration is invalid.
	ErrBadConfig = errors.New("store: configuration is invalid")

	iopsMetrics = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "tigris_gtm",
		Subsystem: "glue",
		Name:      "store_iops",
		Help:      "The number of times each store iop is called",
	}, []string{"driver", "action"})
)

// Interface defines the calls for a generic key value storage interface. This can be
// implemented with an in-memory, on-disk, or in-database storage backend.
type Interface interface {
	// Delete removes a value from the store by key.
	Delete(ctx context.Context, key string) error

	// Exists returns nil if the key exists, ErrNotFound if it does not exist.
	Exists(ctx context.Context, key string) error

	// Get returns the value of a key assuming that value exists and has not expired.
	Get(ctx context.Context, key string) ([]byte, error)

	// Set puts a value into the store that expires according to its expiry.
	Set(ctx context.Context, key string, value []byte) error

	// List lists the keys in this keyspace optionally matching by a prefix.
	List(ctx context.Context, prefix string) ([]string, error)
}

func z[T any]() T { return *new(T) }

type JSON[T any] struct {
	Underlying Interface
	Prefix     string
}

func (j *JSON[T]) Delete(ctx context.Context, key string) error {
	if j.Prefix != "" {
		key = j.Prefix + "/" + key
	}

	return j.Underlying.Delete(ctx, key)
}

func (j *JSON[T]) Exists(ctx context.Context, key string) error {
	if j.Prefix != "" {
		key = j.Prefix + "/" + key
	}

	return j.Underlying.Exists(ctx, key)
}

func (j *JSON[T]) Get(ctx context.Context, key string) (T, error) {
	if j.Prefix != "" {
		key = j.Prefix + "/" + key
	}

	data, err := j.Underlying.Get(ctx, key)
	if err != nil {
		return z[T](), err
	}

	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return z[T](), fmt.Errorf("%w: %w", ErrCantDecode, err)
	}

	return result, nil
}

func (j *JSON[T]) Set(ctx context.Context, key string, value T) error {
	if j.Prefix != "" {
		key = j.Prefix + "/" + key
	}

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCantEncode, err)
	}

	if err := j.Underlying.Set(ctx, key, data); err != nil {
		return err
	}

	return nil
}

func (j *JSON[T]) List(ctx context.Context, prefix string) ([]string, error) {
	fullPrefix := j.Prefix + "/" + prefix
	keys, err := j.Underlying.List(ctx, fullPrefix)
	if err != nil {
		return nil, err
	}

	// Strip the full prefix from each key.
	result := make([]string, 0, len(keys))
	for _, k := range keys {
		result = append(result, strings.TrimPrefix(k, fullPrefix))
	}

	return result, nil
}
