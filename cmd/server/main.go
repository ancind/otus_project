package main

import (
	"flag"
	"github.com/ancind/otus_project/internal"
	"github.com/ancind/otus_project/pkg/image"
	"github.com/ancind/otus_project/pkg/logging"
	lru "github.com/hashicorp/golang-lru"
	"io/ioutil"
	"os"
	"time"
)

var (
	addr            string
	connectTimeout  time.Duration
	requestTimeout  time.Duration
	shutdownTimeout time.Duration
	cacheDir        string
	cacheSize       int
)

func init() {
	flag.StringVar(&addr, "addr", "127.0.0.1:8080", "App addr")
	flag.DurationVar(&connectTimeout, "connect-timeout", 25*time.Second, "Сonnection timeout")
	flag.DurationVar(&requestTimeout, "request-timeout", 25*time.Second, "Request timeout")
	flag.DurationVar(&shutdownTimeout, "shutdown-timeout", 30*time.Second, "Graceful shutdown timeout")
	flag.StringVar(&cacheDir, "cache-dir", "", "Path to Cache dir")
	flag.IntVar(&cacheSize, "cache-size", 5, "Size of cache")
}

func main() {
	flag.Parse()

	logger := logging.DefaultLogger
	imgGetter := image.NewImageGetter(logger, connectTimeout*time.Second, requestTimeout*time.Second)
	resizer := image.NewResizer()

	// 2. If cacheDir isn't provided, then use Temporary Dir
	if cacheDir == "" {
		var err error
		cacheDir, err = ioutil.TempDir("", "")
		if err != nil {
			logger.WithError(err).Fatal(err)
		}
		defer func() {
			if err := os.RemoveAll(cacheDir); err != nil {
				logger.WithError(err).Error("failed to remove cache dir")
			}
		}()
	}

	cache, err := lru.NewWithEvict(cacheSize, func(key interface{}, value interface{}) {
		if path, ok := value.(string); ok {
			defer func() {
				if err := os.Remove(path); err != nil {
					logger.WithError(err).Fatal("failed to remove item from cache")
				}
			}()
		}
	})
	if err != nil {
		logger.WithError(err).Fatal("failed to setup cache")
	}

	server := internal.NewHttp(addr, imgGetter, resizer, cacheDir, cache)
	server.Start()
}
