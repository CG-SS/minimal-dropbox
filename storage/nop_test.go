package storage

import (
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNopStoreDoesNothing(t *testing.T) {
	store := newNopStorage(zerolog.Nop())
	assert.NotNil(t, store)

	filename := "file1.txt"
	fileBytes := []byte("Testing nop")

	err := store.StoreFile(filename, fileBytes)
	assert.NoError(t, err)

	files, err := store.ListFiles()
	assert.NoError(t, err)
	assert.Nil(t, files)

	file, err := store.LoadFile(filename)
	assert.NoError(t, err)
	assert.Nil(t, file)
}
