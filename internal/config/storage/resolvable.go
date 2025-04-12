package storage

import (
	"context"
)

type resolvable struct {
	resolver *Resolver
}

func Resolve(resolver *Resolver) Storage {
	return &resolvable{
		resolver: resolver,
	}
}

func (r *resolvable) Exists(path string) (bool, error) {
	return r.resolver.Resolve(path).Exists(path)
}

func (r *resolvable) Get(ctx context.Context, path string) ([]byte, error) {
	return r.resolver.Resolve(path).Get(ctx, path)
}
