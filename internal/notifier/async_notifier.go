package notifier

import (
	"context"
	"github.com/artarts36/sentry-notifier/internal/metrics"
	"github.com/artarts36/sentry-notifier/internal/sentry"
	"log/slog"
	"time"
)

type AsyncNotifier struct {
	notifier Notifier
	queue    chan sentry.Payload

	running bool
}

func NewAsyncNotifier(notifier Notifier, metr *metrics.Notifier) *AsyncNotifier {
	const queueSize = 500

	queue := make(chan sentry.Payload, queueSize)

	metr.ObserveAsyncQueueSize(func() int {
		return len(queue)
	})

	metr.SetAsyncQueueCapacity(queueSize)

	return &AsyncNotifier{
		notifier: notifier,
		queue:    queue,
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
		n.listen()
	}()
}

func (n *AsyncNotifier) listen() {
	for pl := range n.queue {
		n.notify(pl)
	}
}

func (n *AsyncNotifier) notify(pl sentry.Payload) {
	const timeout = time.Second * 10

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	slog.
		With(slog.Any("payload", pl)).
		DebugContext(ctx, "[async-notifier] processing payload")

	err := n.notifier.Notify(ctx, pl)
	if err != nil {
		slog.
			With(slog.Any("err", err)).
			ErrorContext(ctx, "[async-notifier] failed to listen")
	}

	cancel()
}

