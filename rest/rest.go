package rest

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	"mininal-dropbox/storage"
)

type Server interface {
	Start()
	Stop(ctx context.Context) error
	ErrChan() <-chan error
}

func NewServer(cfg Config, store storage.Storage, logging zerolog.Logger) (Server, error) {
	sys := cfg.System

	switch sys {
	case Gin:
		return newGinServer(cfg, store, logging)
	case Nop:
		return newNopServer(), nil
	default:
		return nil, fmt.Errorf("unknown server system: %s", sys)
	}
}
