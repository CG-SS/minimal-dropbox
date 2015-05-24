package storage

import "fmt"

type memoryStorage struct {
	fileMap map[string][]byte
}

func newMemoryStorage() Storage {
	return memoryStorage{
		fileMap: make(map[string][]byte),
	}
}

func (m memoryStorage) StoreFile(filename string, data []byte) error {
	m.fileMap[filename] = data

	return nil
}

func (m memoryStorage) LoadFile(filename string) ([]byte, error) {
	fileBytes, exists := m.fileMap[filename]
	if !exists {
		return nil, fmt.Errorf("file %s does not exists", filename)
	}

	return fileBytes, nil
}

func (m memoryStorage) ListFiles() ([]string, error) {
	filenames := make([]string, 0, len(m.fileMap))

	for k := range m.fileMap {
		filenames = append(filenames, k)
	}

	return filenames, nil
}
