package storage

import (
	"fmt"
)

type Storage interface {
	StoreFile(filename string, data []byte) error
	LoadFile(filename string) ([]byte, error)
	ListFiles() ([]string, error)
}

func NewStorage(cfg Config) (Storage, error) {
	switch cfg.System {
	case FileSystem:
		return newFileSystemStorage(cfg)
	case Memory:
		return newMemoryStorage(), nil
	case Nop:
		return newNopStorage(), nil
	default:
		return nil, fmt.Errorf("unkown storage system %s", cfg.System)
	}
}
