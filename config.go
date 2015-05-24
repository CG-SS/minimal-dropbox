package main

import (
	"github.com/kelseyhightower/envconfig"
	"log"
	"mininal-dropbox/rest"
	"mininal-dropbox/storage"
)

type Config struct {
	Store          storage.Config
	Rest           rest.Config 
	LoggingEnabled bool `envconfig:"LOGGING_ENABLED" default:"true"`
}

func NewConfig() Config {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("could not process envconfig: %v", err)
	}

	return cfg
}
