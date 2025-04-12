package storage

import (
	"context"
	"errors"
	"fmt"
	"os"
)

type Filesystem struct {
}

func NewFilesystem() *Filesystem {
	return &Filesystem{}
}

func (f *Filesystem) Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, fmt.Errorf("stat: %w", err)
	}
	return true, nil
}

func (f *Filesystem) Get(_ context.Context, key string) ([]byte, error) {
	return os.ReadFile(key)
}
