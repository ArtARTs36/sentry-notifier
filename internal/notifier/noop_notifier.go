package notifier

import (
	"context"
	"log/slog"

	"github.com/artarts36/sentry-notifier/internal/sentry"
)

type NoopNotifier struct {
}

func NewNoopNotifier() *NoopNotifier {
	return &NoopNotifier{}
}

func (n *NoopNotifier) Notify(ctx context.Context, pl sentry.Payload) error {
	slog.
		With(slog.String("payload_id", pl.GetID())).
		DebugContext(ctx, "[noop-notifier] skip payload")

	return nil
}

func (n *NoopNotifier) Close() {
}
