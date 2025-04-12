package storage

import (
	"context"
	"fmt"
	"os"
	"strings"
)

type Env struct{}

func NewEnv() *Env {
	return &Env{}
}

func (h *Env) Exists(path string) (bool, error) {
	_, exists := os.LookupEnv(h.key(path))
	return exists, nil
}

func (h *Env) Get(_ context.Context, path string) ([]byte, error) {
	varName := h.key(path)

	val, ok := os.LookupEnv(varName)
	if !ok {
		return nil, fmt.Errorf("environment variable %q not found", varName)
	}

	return []byte(val), nil
}

func (h *Env) key(path string) string {
	varName := strings.TrimPrefix(path, "env://")
	varName = strings.TrimPrefix(varName, "$")

	return varName
}
