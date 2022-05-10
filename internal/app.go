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
		params := r.URL.Query()
		url := params.Get("url")
		rawWidth, rawHeight := params.Get("width"), params.Get("height")

		width, err := strconv.Atoi(rawWidth)
		if err != nil {

			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		height, err := strconv.Atoi(rawHeight)
		if err != nil {

			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		img, _ := a.handle(r.Context(), url, r.Header, width, height)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		w.Header().Add("Content-Type", "image/jpeg")
		w.Header().Set("Content-Length", strconv.Itoa(len(img)))

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

	// 2. If not found in Cache, then try to fetch Image
	img, err := a.imgGetter.Get(ctx, url, header)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch image")
	}

	// 2. Transform Image
	img, err = a.resizer.Resize(img, width, height)
	if err != nil {
		return nil, errors.Wrap(err, "failed to crop image")
	}

	// 3. Save transformed Image
	imgPath := path.Join(a.cacheDir, cacheKey+".jpeg")
	err = ioutil.WriteFile(imgPath, img, 0600)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save image")
	}

	// 4. Add Image path to Cache
	a.cache.Add(cacheKey, imgPath)

	return img, nil
}
