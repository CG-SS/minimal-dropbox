package main

import (
	"context"
	"log"
	"mininal-dropbox/rest"
	"mininal-dropbox/storage"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := NewConfig()

	if !cfg.LoggingEnabled {
		log.SetOutput(nil)
	}

	store, err := storage.NewStorage(cfg.Store)
	if err != nil {
		log.Fatal("failed initializing database")
	}

	log.Printf("starting server at %s:%d", cfg.Rest.Host, cfg.Rest.Port)

	restServer, err := rest.NewServer(cfg.Rest, store)
	if err != nil {
		log.Fatal("failed initializing server")
	}

	go restServer.Start()

	sigs := make(chan os.Signal, 1)
	defer close(sigs)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case serverErr := <-restServer.ErrChan():
		log.Printf("got error from server, shutting down: %v", serverErr)
	case s := <-sigs:
		log.Printf("received signal, shutting down: %v", s)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := restServer.Stop(ctx)
		if err != nil {
			log.Fatalf("got error while trying to shutdown http server: %v", err)
		}
	}
}
