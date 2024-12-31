package storage

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStorageListFiles(t *testing.T) {
	files := []string{
		"test1.txt",
		"test2.txt",
		"test3.txt",
		"test4.txt",
		"test5.txt",
	}
	fileMap := make(map[string][]byte, len(files))

	for _, file := range files {
		fileMap[file] = []byte{}
	}

	store := memoryStorage{fileMap: fileMap}

	listFiles, err := store.ListFiles()
	assert.NoError(t, err)
	for _, file := range listFiles {
		assert.Contains(t, file, file)
	}
}

func TestMemoryStorageStoreFile(t *testing.T) {
	fileMap := make(map[string][]byte, 1)

	store := memoryStorage{fileMap: fileMap}
	filename := "test1.txt"
	err := store.StoreFile("test1.txt", bytes.NewBufferString("test"))
	assert.NoError(t, err)

	_, exists := fileMap[filename]
	assert.True(t, exists)
}

func TestMemoryStorageLoadFile(t *testing.T) {
	filename := "test1.txt"
	value := []byte("Hello there")

	fileMap := make(map[string][]byte, 1)
	fileMap[filename] = value

	store := memoryStorage{fileMap: fileMap}
	file, err := store.LoadFile(filename)

	all, err := io.ReadAll(file)

	assert.NoError(t, err)
	assert.Equal(t, value, all)
}
