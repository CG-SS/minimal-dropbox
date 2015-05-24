package storage

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestEnvconfigLoadsEnvVar(t *testing.T) {
	var cfg Config

	storageSysVal := Nop
	managedDirVal := "./testing"

	err := os.Setenv("STORAGE_SYSTEM", string(storageSysVal))
	assert.NoError(t, err)
	err = os.Setenv("STORAGE_MANAGED_DIR", managedDirVal)
	assert.NoError(t, err)

	err = envconfig.Process("", &cfg)
	assert.NoError(t, err)

	assert.Equal(t, storageSysVal, cfg.System)
	assert.Equal(t, managedDirVal, cfg.ManagedDir)
}
