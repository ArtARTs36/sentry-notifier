package injector

import (
	"github.com/artarts36/sentry-notifier/internal/config/cfg"
	"github.com/artarts36/sentry-notifier/internal/notifier"
)

type NotifyDefaultStrategy struct {
}

func NewNotifyDefaultStrategy() *NotifyDefaultStrategy {
	return &NotifyDefaultStrategy{}
}

func (*NotifyDefaultStrategy) Inject(config cfg.Config) (cfg.Config, error) {
	if config.Notify.Strategy == "" {
		config.Notify.Strategy = notifier.StrategyDefault
	}

	return config, nil
}
