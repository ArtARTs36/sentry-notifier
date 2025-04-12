package loader

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/artarts36/sentry-notifier/internal/config/cfg"
	"github.com/artarts36/sentry-notifier/internal/config/injector"
	"github.com/artarts36/sentry-notifier/internal/config/parser"
	"github.com/artarts36/sentry-notifier/internal/config/storage"
)

type Loader struct {
	storage  storage.Storage
	parsers  map[string]parser.Parser
	injector injector.Injector
}

func New(storage storage.Storage, parsers map[string]parser.Parser, injector injector.Injector) *Loader {
	return &Loader{
		storage:  storage,
		parsers:  parsers,
		injector: injector,
	}
}

func (l *Loader) Load(ctx context.Context, path string) (cfg.Config, error) {
	content, err := l.storage.Get(ctx, path)
	if err != nil {
		return cfg.Config{}, fmt.Errorf("failed to get config from storage: %w", err)
	}

	pars, err := l.getParser(path)
	if err != nil {
		return cfg.Config{}, err
	}

	config, err := pars.Parse(ctx, content)
	if err != nil {
		return cfg.Config{}, fmt.Errorf("failed to parse config: %w", err)
	}

	if err = config.Validate(); err != nil {
		return cfg.Config{}, fmt.Errorf("config invalid: %w", err)
	}

	config, err = l.injector.Inject(config)
	if err != nil {
		return cfg.Config{}, fmt.Errorf("failed to inject config: %w", err)
	}

	return config, nil
}

func (l *Loader) getParser(path string) (parser.Parser, error) {
	ext := filepath.Ext(path)[1:]
	if ext == "" {
		return nil, fmt.Errorf("invalid config file extension")
	}

	p, ok := l.parsers[ext]
	if !ok {
		return nil, fmt.Errorf("config with extension %q not supported", ext)
	}

	return p, nil
}
