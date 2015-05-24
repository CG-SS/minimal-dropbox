package storage

import (
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewStorageFailsUnknown(t *testing.T) {
	var cfg Config
	cfg.System = "test"

	storage, err := NewStorage(cfg, zerolog.Nop())
	assert.Error(t, err)
	assert.Nil(t, storage)
}
