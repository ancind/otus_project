package main

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/ancind/otus_project/pkg/app"
	lru "github.com/hashicorp/golang-lru"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	addr            string
	connectTimeout  time.Duration
	requestTimeout  time.Duration
	shutdownTimeout time.Duration
	cacheSize       int
)

func init() {
	flag.StringVar(&addr, "addr", "0.0.0.0:8080", "App addr")
	flag.DurationVar(&connectTimeout, "connect-timeout", 25*time.Second, "Connection timeout")
	flag.DurationVar(&requestTimeout, "request-timeout", 25*time.Second, "Request timeout")
	flag.DurationVar(&shutdownTimeout, "shutdown-timeout", 30*time.Second, "Graceful shutdown timeout")
	flag.IntVar(&cacheSize, "cache-size", 5, "Size of cache")
}

var shaCommit = "local"

func main() {
	flag.Parse()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := log.With().Str("sha_commit", shaCommit).Logger()

	cache, err := lru.NewWithEvict(cacheSize, func(key interface{}, value interface{}) {
		if path, ok := value.(string); ok {
			defer func() {
				if err := os.Remove(path); err != nil {
					logger.Fatal().Err(err).Msg("failed to remove item from cache")
				}
			}()
		}
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to setup cache")
	}

	srv, err := app.NewServer(logger, connectTimeout, requestTimeout, cache)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to start server")
	}

	ctx := log.Logger.WithContext(context.Background())
	if err := srv.Listen(ctx); err != nil {
		logger.Fatal().Err(err).Msg("failed listening port")
	}
}
