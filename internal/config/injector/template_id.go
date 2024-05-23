package injector

import (
	"fmt"
	"github.com/artarts36/sentry-notifier/internal/config/cfg"
)

type TemplateID struct {
}

func NewTemplateID() *TemplateID {
	return &TemplateID{}
}

func (*TemplateID) Inject(config cfg.Config) (cfg.Config, error) {
	i := 0

	for templateIndex, messages := range config.Notify.On {
		for messageIndex := range messages {
			config.Notify.On[templateIndex][messageIndex].MessageTemplateID = fmt.Sprintf("%d", i)

			i++
		}
	}

	return config, nil
}
