package notifier

import (
	"github.com/artarts36/sentry-notifier/internal/messenger"
	"github.com/artarts36/sentry-notifier/internal/metrics"
	"github.com/artarts36/sentry-notifier/internal/template"
)

func CreateNotifier(
	messengers map[string][]messenger.Messenger,
	renderer *template.Renderer,
	cfg Config,
	metr *metrics.Group,
) Notifier {
	switch cfg.Strategy {
	case StrategyImmediately:
		return NewImmediatelyNotifier(messengers, renderer, cfg, metr.Messages)
	case StrategyAsync:
		return NewAsyncNotifier(NewImmediatelyNotifier(messengers, renderer, cfg, metr.Messages), metr.Notifier)
	case StrategyNull:
		return NewNoopNotifier()
	}

	return NewImmediatelyNotifier(messengers, renderer, cfg, metr.Messages)
}
