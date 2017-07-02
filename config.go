package main

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"log"
	"mininal-dropbox/rest"
	"mininal-dropbox/storage"
)

type ZerologLevel int8

func (zl *ZerologLevel) Decode(value string) error {
	level, err := zerolog.ParseLevel(value)
	if err != nil {
		return fmt.Errorf("failed parsing zerolog level: %w", err)
	}

	*zl = ZerologLevel(level)

	return nil
}

type Config struct {
	Store          storage.Config
	Rest           rest.Config
	LoggingEnabled bool         `envconfig:"LOGGING_ENABLED" default:"true"`
	LoggingLevel   ZerologLevel `envconfig:"LOGGING_LEVEL" default:"debug"`
}

func NewConfig() Config {
	var cfg Config

	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("could not process envconfig: %v", err)
	}

	return cfg
}
