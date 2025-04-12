package storage

import "context"

type Storage interface {
	Exists(path string) (bool, error)
	Get(ctx context.Context, key string) ([]byte, error)
}
