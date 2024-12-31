package storage

import (
	"bytes"
	"io"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestNopStoreDoesNothing(t *testing.T) {
	store := newNopStorage(zerolog.Nop())
	assert.NotNil(t, store)

	filename := "file1.txt"
	fileBytes := []byte("Testing nop")

	err := store.StoreFile(filename, bytes.NewReader(fileBytes))
	assert.NoError(t, err)

	files, err := store.ListFiles()
	assert.NoError(t, err)
	assert.Nil(t, files)

	file, err := store.LoadFile(filename)
	assert.NoError(t, err)

	fileAll, err := io.ReadAll(file)
	assert.NoError(t, err)

	assert.NotEqual(t, fileAll, []byte("Testing nop"))
	err = file.Close()
	assert.NoError(t, err)
}
