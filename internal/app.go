package internal

import (
	"context"
	"fmt"
	"github.com/ancind/otus_project/pkg/image"
	"github.com/ancind/otus_project/pkg/util"
	"github.com/gorilla/mux"
	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type app struct {
	imgGetter image.Getter
	resizer   *image.Resizer
	cacheDir  string
	cache     *lru.Cache
}

func NewApp(ig image.Getter, r *image.Resizer, cacheDir string, cache *lru.Cache) *app {
	return &app{imgGetter: ig, resizer: r, cacheDir: cacheDir, cache: cache}
}

func (a *app) Run() http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var pu string
		parsedURL, err := url.Parse(vars["imageUrl"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		pu = parsedURL.String()
		if parsedURL.Scheme == "" {
			pu = fmt.Sprintf("http://%s", parsedURL.String())
		}

		width, err := strconv.Atoi(vars["width"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		height, err := strconv.Atoi(vars["height"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		img, err := a.handle(r.Context(), pu, r.Header, width, height)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		w.Header().Add("Content-Type", "image/jpeg")
		w.Header().Set("Content-Length", strconv.Itoa(len(img)))

		if _, err := w.Write(img); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
		}
	}
	router := mux.NewRouter()
	router.HandleFunc("/fill/{width:[0-9]+}/{height:[0-9]+}/{imageUrl:.*}", f)
	http.Handle("/", router)

	return router
}

func (a *app) handle(ctx context.Context, url string, header http.Header, width, height int) ([]byte, error) {
	// 1. Try to find image in Cache
	s := fmt.Sprintf("%s|%d|%d", url, width, height)
	cacheKey := util.GetHash(s)

	if imgPath, found := a.cache.Get(cacheKey); found {
		img, err := ioutil.ReadFile(imgPath.(string))
		return img, err
	}

	img, err := a.imgGetter.Get(ctx, url, header)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get image")
	}

	img, err = a.resizer.Resize(img, width, height)
	if err != nil {
		return nil, errors.Wrap(err, "failed to transform image")
	}

	imgPath := path.Join(a.cacheDir, cacheKey+".jpeg")
	err = ioutil.WriteFile(imgPath, img, 0600)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save image")
	}

	// 4. Add Image path to Cache
	a.cache.Add(cacheKey, imgPath)

	return img, nil
}
