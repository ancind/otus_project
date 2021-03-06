package image

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
)

type Service interface {
	ResizeImage(context.Context, string, http.Header, int, int) ([]byte, error)
}

type service struct {
	imgGetter Getter
	resizer   Transformer
	cache     *lru.Cache
}

func NewService(ig Getter, t Transformer, cache *lru.Cache) Service {
	return &service{imgGetter: ig, resizer: t, cache: cache}
}

func (s *service) ResizeImage(ctx context.Context, url string, header http.Header, width, height int) ([]byte, error) {
	// 1. Try to find image in Cache
	cacheKey := fmt.Sprintf("%s|%d|%d", url, width, height)

	if imgPath, found := s.cache.Get(cacheKey); found {
		img, err := ioutil.ReadFile(imgPath.(string))
		return img, err
	}

	img, err := s.imgGetter.Get(ctx, url, header)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get image")
	}

	img, err = s.resizer.Resize(img, width, height)
	if err != nil {
		return nil, errors.Wrap(err, "failed to transform image")
	}

	// 4. Add Image path to Cache
	s.cache.Add(cacheKey, img)

	return img, nil
}
