package main

import (
	"github.com/ancind/otus_project/internal"
	"github.com/ancind/otus_project/pkg/image"
	"github.com/ancind/otus_project/pkg/logging"
	lru "github.com/hashicorp/golang-lru"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"time"
)

func main() {
	// 1. Setup required Units
	logger := logging.DefaultLogger

	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		logger.WithError(err).Error("failed to read .env")
	}

	connectTimeout := viper.Get("CONNECT_TIMEOUT").(time.Duration)
	requestTimeout := viper.Get("REQUEST_TIMEOUT").(time.Duration)
	cacheDir := viper.Get("CACHE_DIR").(string)
	cacheSize := viper.Get("CACHE_SIZE").(int)

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

	// 3. Setup Cache
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

	server := internal.NewHttp(imgGetter, resizer, cacheDir, cache)
	server.Start()
}
