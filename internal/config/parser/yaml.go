package parser

import (
	"context"
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/artarts36/sentry-notifier/internal/config/cfg"
)

type YAML struct {
}

func NewYAML() *YAML {
	return &YAML{}
}

func (l *YAML) Parse(_ context.Context, content []byte) (cfg.Config, error) {
	var conf cfg.Config

	err := yaml.Unmarshal(content, &conf)
	if err != nil {
		return conf, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return conf, nil
}
