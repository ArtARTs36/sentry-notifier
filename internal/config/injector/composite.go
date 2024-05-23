package injector

import "github.com/artarts36/sentry-notifier/internal/config/cfg"

type Composite struct {
	injectors []Injector
}

func NewComposite(injectors []Injector) *Composite {
	return &Composite{
		injectors: injectors,
	}
}

func (c *Composite) Inject(config cfg.Config) (cfg.Config, error) {
	conf := config

	for _, injector := range c.injectors {
		var err error
		conf, err = injector.Inject(conf)
		if err != nil {
			return cfg.Config{}, err
		}
	}

	return conf, nil
}
