package notifier

import (
	"context"
	"github.com/artarts36/sentry-notifier/internal/sentry"
	"log/slog"
	"time"
)

type AsyncNotifier struct {
	notifier Notifier
	queue    chan sentry.Payload

	running bool
}

func NewAsyncNotifier(notifier Notifier) *AsyncNotifier {
	const queueSize = 500

	return &AsyncNotifier{
		notifier: notifier,
		queue:    make(chan sentry.Payload, queueSize),
	}
}

func (n *AsyncNotifier) Notify(ctx context.Context, payload sentry.Payload) error {
	n.run()

	slog.
		With(slog.Any("payload", payload)).
		DebugContext(ctx, "[async-notifier] push payload to queue")

	n.queue <- payload

	return nil
}

func (n *AsyncNotifier) Close() {
	close(n.queue)
}

func (n *AsyncNotifier) run() {
	if n.running {
		return
	}

	n.running = true

	go func() {
		n.notify()
	}()
}

func (n *AsyncNotifier) notify() {
	const timeout = time.Second * 10

	for { //nolint: gosimple // not need
		select {
		case pl := <-n.queue:
			ctx, cancel := context.WithTimeout(context.Background(), timeout)

			slog.
				With(slog.Any("payload", pl)).
				DebugContext(ctx, "[async-notifier] processing payload")

			err := n.notifier.Notify(ctx, pl)
			if err != nil {
				slog.
					With(slog.String("err", err.Error())).
					ErrorContext(ctx, "[async-notifier] failed to notify")
			}

			cancel()
		}
	}
}
