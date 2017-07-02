package storage

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"path"
	"path/filepath"
)

type fileSystem struct {
	managedDir string
	logging    zerolog.Logger
}

func newFileSystemStorage(cfg Config, logging zerolog.Logger) (Storage, error) {
	err := validateDirPath(cfg.ManagedDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create file system storage: %w", err)
	}

	return fileSystem{
		managedDir: cfg.ManagedDir,
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

func (f fileSystem) LoadFile(filename string) ([]byte, error) {
	fileBytes, err := os.ReadFile(path.Join(f.managedDir, filename))
	if err != nil {
		return nil, fmt.Errorf("failed opening file: %w", err)
	}

	f.logging.Debug().Msg(fmt.Sprintf("loaded file: %s", filename))

	return fileBytes, nil
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

func (f fileSystem) StoreFile(filename string, data []byte) error {
	err := os.WriteFile(filepath.Join(f.managedDir, filename), data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	f.logging.Debug().Msg(fmt.Sprintf("wrote file: %s", filename))

	return nil
}
