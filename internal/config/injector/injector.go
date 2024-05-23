package injector

import "github.com/artarts36/sentry-notifier/internal/config/cfg"

type Injector interface {
	Inject(config cfg.Config) (cfg.Config, error)
}
