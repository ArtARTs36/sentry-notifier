package notifier

import (
	"github.com/artarts36/sentry-notifier/internal/messenger"
	"github.com/artarts36/sentry-notifier/internal/template"
)

func CreateNotifier(
	messengers map[string][]messenger.Messenger,
	renderer *template.Renderer,
	cfg Config,
) Notifier {
	switch cfg.Strategy {
	case StrategyImmediately:
		return NewImmediatelyNotifier(messengers, renderer, cfg)
	case StrategyAsync:
		return NewAsyncNotifier(NewImmediatelyNotifier(messengers, renderer, cfg))
	}

	return NewImmediatelyNotifier(messengers, renderer, cfg)
}
