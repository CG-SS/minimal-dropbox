package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	"mininal-dropbox/rest"
	"mininal-dropbox/storage"
)

func main() {
	cfg := NewConfig()

	logging := zerolog.Nop()
	if cfg.LoggingEnabled {
		zerolog.SetGlobalLevel(zerolog.Level(cfg.LoggingLevel))

		logging = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	store, err := storage.NewStorage(cfg.Store, logging)
	if err != nil {
		logging.Fatal().Err(err).Msg("failed initializing database")
	}

	logging.Info().Msg(fmt.Sprintf("starting server at %s:%d", cfg.Rest.Host, cfg.Rest.Port))

	restServer, err := rest.NewServer(cfg.Rest, store, logging)
	if err != nil {
		logging.Fatal().Err(err).Msg("failed initializing server")
	}

	go restServer.Start()

	sigs := make(chan os.Signal, 1)
	defer close(sigs)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case serverErr := <-restServer.ErrChan():
		logging.Error().Err(serverErr).Msg("got error from server, shutting down...")
	case s := <-sigs:
		logging.Info().Msg(fmt.Sprintf("received signal %v, shutting down...", s))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := restServer.Stop(ctx)
		if err != nil {
			logging.Fatal().Err(err).Msg("got error while trying to shutdown http server")
		}
	}
}
