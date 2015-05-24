package rest

import (
	"context"
	"mininal-dropbox/storage"
)

type Server interface {
	Start()
	Stop(ctx context.Context) error
	ErrChan() <-chan error
}

func NewServer(cfg Config, store storage.Storage) (Server, error) {
	if cfg.System == Gin {
		return newGinServer(cfg, store)
	}

	return newNopServer(), nil
}
