package parser

import (
	"context"

	"github.com/artarts36/sentry-notifier/internal/config/cfg"
)

type Parser interface {
	Parse(ctx context.Context, content []byte) (cfg.Config, error)
}
