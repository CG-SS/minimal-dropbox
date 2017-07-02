package storage

import (
	"fmt"
	"github.com/rs/zerolog"
)

type Storage interface {
	DeleteFile(filename string) error
	StoreFile(filename string, data []byte) error
	LoadFile(filename string) ([]byte, error)
	ListFiles() ([]string, error)
}

func NewStorage(cfg Config, logging zerolog.Logger) (Storage, error) {
	switch cfg.System {
	case FileSystem:
		return newFileSystemStorage(cfg, logging)
	case Memory:
		return newMemoryStorage(), nil
	case Nop:
		return newNopStorage(logging), nil
	default:
		return nil, fmt.Errorf("unkown storage system %s", cfg.System)
	}
}
