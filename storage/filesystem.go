package storage

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/rs/zerolog"
)

type fileSystem struct {
	managedDir string
	logging    zerolog.Logger
	bufferSize int
}

func newFileSystemStorage(cfg Config, logging zerolog.Logger) (Storage, error) {
	err := validateDirPath(cfg.ManagedDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create file system storage: %w", err)
	}

	return fileSystem{
		managedDir: cfg.ManagedDir,
		bufferSize: cfg.BufferSize,
		logging:    logging,
	}, nil
}

func validateDirPath(dir string) error {
	stat, err := os.Stat(dir)
	if err != nil {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failing creating dir: %w", err)
		}

		return nil
	}

	if !stat.IsDir() {
		return fmt.Errorf("provided path is not a dir")
	}

	return nil
}

func (f fileSystem) DeleteFile(filename string) error {
	err := os.Remove(path.Join(f.managedDir, filename))
	if err != nil {
		return fmt.Errorf("could not delete file %s: %w", filename, err)
	}

	return nil
}

func (f fileSystem) LoadFile(filename string) (io.ReadCloser, error) {
	filePath := path.Join(f.managedDir, filename)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed opening file: %w", err)
	}

	f.logging.Debug().Msg(fmt.Sprintf("loaded file: %s", filename))

	return file, nil
}

func (f fileSystem) ListFiles() ([]string, error) {
	openManagedDir, err := os.Open(f.managedDir)
	if err != nil {
		return nil, fmt.Errorf("failed opening managed dir: %w", err)
	}

	filenames, err := openManagedDir.Readdirnames(-1)
	if err != nil {
		return nil, fmt.Errorf("failed listing files: %w", err)
	}

	f.logging.Debug().Msg(fmt.Sprintf("found %d files", len(filenames)))

	return filenames, nil
}

func (f fileSystem) StoreFile(filename string, reader io.Reader) error {
	filePath := path.Join(f.managedDir, filename)

	newFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	buffer := make([]byte, f.bufferSize)
	for {
		_, err := reader.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				return fmt.Errorf("failed to read chunk for filename '%s': %w", filename, err)
			}
		}

		_, err = newFile.Write(buffer)
		if err != nil {
			return fmt.Errorf("failed to write file chunk for filename '%s': %w", filename, err)
		}
	}

	f.logging.Debug().Msg(fmt.Sprintf("wrote file: %s", filename))

	return nil
}
