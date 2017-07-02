package rest

import (
	"context"
	"net/http"
)

type nopServer struct {
	errChan chan error
	stopper chan struct{}
}

func newNopServer() Server {
	return &nopServer{
		errChan: make(chan error),
		stopper: make(chan struct{}),
	}
}

func (s nopServer) Start() {
	<-s.stopper
	s.errChan <- http.ErrServerClosed

	close(s.errChan)
}

func (s nopServer) Stop(_ context.Context) error {
	close(s.stopper)

	return nil
}

func (s nopServer) ErrChan() <-chan error {
	return s.errChan
}
