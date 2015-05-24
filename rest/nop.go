package rest

import "context"

type nopServer struct {
	errChan chan error
}

func newNopServer() Server {
	return &nopServer{
		errChan: make(chan error),
	}
}

func (s *nopServer) Start() {}

func (s *nopServer) Stop(_ context.Context) error {
	return nil
}

func (s *nopServer) ErrChan() <-chan error {
	return s.errChan
}
