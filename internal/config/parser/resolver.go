package parser

import (
	"errors"
	"fmt"
	"path/filepath"
)

type Resolver struct {
	json *JSON
	yaml *YAML
}

func NewResolver() *Resolver {
	return &Resolver{
		json: NewJSON(),
		yaml: NewYAML(),
	}
}

func (r *Resolver) Resolve(path string, content []byte) (Parser, error) {
	ext := filepath.Ext(path)
	if ext != "" {
		if ext == ".json" {
			return r.json, nil
		}
		if ext == ".yaml" {
			return r.yaml, nil
		}
		if ext == ".yml" {
			return r.yaml, nil
		}

		return nil, fmt.Errorf("extension %q unsupported", ext[1:])
	}

	if content[0] == '{' {
		return r.json, nil
	}

	return nil, errors.New("cannot determine file format")
}
