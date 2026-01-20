package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"log/slog"

	"github.com/facebookgo/flagenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/tigrisdata-community/glue/internal/store"
)

var (
	storeBucket = flag.String("store-bucket", "", "The Tigris bucket used to store data")
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
		"store-bucket", *storeBucket,
		"args", flag.Args(),
	)

	if err := run(ctx); err != nil {
		log.Fatal("error:", err)
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
	_ = seenURLs

	return nil
}
