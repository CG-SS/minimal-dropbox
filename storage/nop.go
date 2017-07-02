package storage

import (
	"fmt"
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

func (n nopStorage) LoadFile(filename string) ([]byte, error) {
	n.logging.Debug().Msg(fmt.Sprintf("would load: %s", filename))

	return nil, nil
}

func (n nopStorage) ListFiles() ([]string, error) {
	n.logging.Debug().Msg("would list files")

	return nil, nil
}

func (n nopStorage) StoreFile(filename string, _ []byte) error {
	n.logging.Debug().Msg(fmt.Sprintf("would store: %s", filename))

	return nil
}
