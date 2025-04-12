package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/artarts36/sentry-notifier/internal/config/cfg"
)

type JSON struct {
}

func NewJSON() *JSON {
	return &JSON{}
}

func (l *JSON) Parse(_ context.Context, content []byte) (cfg.Config, error) {
	var conf cfg.Config

	err := json.Unmarshal(content, &conf)
	if err != nil {
		return conf, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return conf, nil
}
