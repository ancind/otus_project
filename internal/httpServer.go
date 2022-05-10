package internal

import (
	"github.com/ancind/otus_project/pkg/image"
	lru "github.com/hashicorp/golang-lru"
	"net/http"
	"time"
)

type HttpServer struct {
	server *http.Server
}

func NewHttp(ig *image.HttpGetter, ir *image.Resizer, cd string, c *lru.Cache) *HttpServer {
	app := NewApp(ig, ir, cd, c)

	httpSrv := &http.Server{
		Addr:         ":8080",
		Handler:      app.Run(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return &HttpServer{httpSrv}
}

func (s *HttpServer) Start() {
	s.server.ListenAndServe()
}
