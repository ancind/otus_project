package app

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/ancind/otus_project/pkg/image"
	"github.com/gorilla/mux"
	"github.com/hashicorp/golang-lru"
	"github.com/rs/zerolog"
)

type Server struct {
	svc    image.Service
	logger zerolog.Logger
	router *mux.Router
}

func NewServer(logger zerolog.Logger, ct, cr time.Duration, cd string, c *lru.Cache) (*Server, error) {
	svc := image.NewService(image.NewImageGetter(logger, ct, cr), image.NewResizer(), cd, c)

	srv := &Server{
		svc:    svc,
		logger: logger,
	}

	srv.createRoute()

	return srv, nil
}

func (s *Server) Listen(ctx context.Context) error {
	httpSrv := &http.Server{
		Addr:         "0.0.0.0:8080",
		Handler:      s.router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	return httpSrv.ListenAndServe()
}

func (s *Server) createRoute() {
	r := mux.NewRouter()
	handlers := NewHandler(s.logger, s.svc)

	r.HandleFunc("/fill/{width:[0-9]+}/{height:[0-9]+}/{imageUrl:.*}", handlers.ImageHandler)

	s.router = r
}
