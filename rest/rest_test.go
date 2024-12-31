package rest

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"mininal-dropbox/storage"
)

func TestNewServerCreatesNopServer(t *testing.T) {
	var cfg Config
	cfg.System = Nop

	nopLogging := zerolog.Nop()

	var storeCfg storage.Config
	storeCfg.System = storage.Nop
	store, err := storage.NewStorage(storeCfg, nopLogging)
	assert.NoError(t, err)

	server, err := NewServer(cfg, store, nopLogging)
	assert.NoError(t, err)

	var nServer *nopServer
	assert.Equal(t, reflect.TypeOf(server).Kind(), reflect.TypeOf(nServer).Kind())
}

func TestNewServerCreatesGinServer(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	var cfg Config
	cfg.System = Gin
	cfg.Port = 12345
	cfg.Host = "localhost"

	nopLogging := zerolog.Nop()

	var storeCfg storage.Config
	storeCfg.System = storage.Nop
	store, err := storage.NewStorage(storeCfg, nopLogging)
	assert.NoError(t, err)

	server, err := NewServer(cfg, store, nopLogging)
	assert.NoError(t, err)

	var nServer *ginServer
	assert.Equal(t, reflect.TypeOf(server).Kind(), reflect.TypeOf(nServer).Kind())
}

func TestNewServerFailsWithUnknownSystem(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	var cfg Config
	cfg.System = "test"
	cfg.Port = 12345
	cfg.Host = "localhost"

	nopLogging := zerolog.Nop()

	var storeCfg storage.Config
	storeCfg.System = storage.Nop
	store, err := storage.NewStorage(storeCfg, nopLogging)
	assert.NoError(t, err)

	server, err := NewServer(cfg, store, nopLogging)
	assert.Error(t, err)
	assert.Nil(t, server)
}
