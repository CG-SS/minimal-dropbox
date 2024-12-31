package storage

import (
	"bytes"
	"fmt"
	"io"
)

type memoryFile struct {
	reader io.Reader
}

func (m memoryFile) Read(p []byte) (n int, err error) {
	return m.reader.Read(p)
}

func (m memoryFile) Close() error {
	return nil
}

type memoryStorage struct {
	fileMap map[string][]byte
}

func newMemoryStorage() Storage {
	return memoryStorage{
		fileMap: make(map[string][]byte),
	}
}

func (m memoryStorage) DeleteFile(filename string) error {
	delete(m.fileMap, filename)

	return nil
}

func (m memoryStorage) StoreFile(filename string, reader io.Reader) error {
	fileBytes, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to store '%s': %w", filename, err)
	}

	m.fileMap[filename] = fileBytes

	return nil
}

func (m memoryStorage) LoadFile(filename string) (io.ReadCloser, error) {
	fileBytes, exists := m.fileMap[filename]
	if !exists {
		return nil, fmt.Errorf("file %s does not exists", filename)
	}

	return memoryFile{bytes.NewReader(fileBytes)}, nil
}

func (m memoryStorage) ListFiles() ([]string, error) {
	filenames := make([]string, 0, len(m.fileMap))

	for k := range m.fileMap {
		filenames = append(filenames, k)
	}

	return filenames, nil
}
