package app

import (
	"github.com/ancind/otus_project/pkg/image"
	"github.com/rs/zerolog"
)

type Handlers struct {
	logger zerolog.Logger
	svc    image.Service
}

func NewHandler(
	logger zerolog.Logger,
	svc image.Service,
) *Handlers {
	return &Handlers{logger: logger, svc: svc}
}
