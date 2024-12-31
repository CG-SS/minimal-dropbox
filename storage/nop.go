package storage

import (
	"bytes"
	"fmt"
	"io"

	"github.com/rs/zerolog"
)

type nopStorage struct {
	logging zerolog.Logger
}

func newNopStorage(logging zerolog.Logger) Storage {
	return nopStorage{
		logging: logging,
	}
}

func (n nopStorage) DeleteFile(filename string) error {
	n.logging.Debug().Msg(fmt.Sprintf("would delete: %s", filename))

	return nil
}

func (n nopStorage) LoadFile(filename string) (io.ReadCloser, error) {
	n.logging.Debug().Msg(fmt.Sprintf("would load: %s", filename))

	return memoryFile{bytes.NewReader([]byte(""))}, nil
}

func (n nopStorage) ListFiles() ([]string, error) {
	n.logging.Debug().Msg("would list files")

	return nil, nil
}

func (n nopStorage) StoreFile(filename string, _ io.Reader) error {
	n.logging.Debug().Msg(fmt.Sprintf("would store: %s", filename))

	return nil
}
