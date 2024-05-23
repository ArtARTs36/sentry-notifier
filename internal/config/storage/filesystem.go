package storage

import (
	"context"
	"os"
)

type Filesystem struct {
}

func NewFilesystem() *Filesystem {
	return &Filesystem{}
}

func (f *Filesystem) Get(_ context.Context, key string) ([]byte, error) {
	return os.ReadFile(key)
}
