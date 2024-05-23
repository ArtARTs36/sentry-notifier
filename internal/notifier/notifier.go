package notifier

import (
	"context"

	"github.com/artarts36/sentry-notifier/internal/sentry"
)

type Notifier interface {
	Notify(ctx context.Context, payload sentry.Payload) error
	Close()
}
