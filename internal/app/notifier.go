package app

import (
	"github.com/artarts36/sentry-notifier/internal/config/cfg"
	"github.com/artarts36/sentry-notifier/internal/messenger"
	"github.com/artarts36/sentry-notifier/internal/notifier"
	"github.com/artarts36/sentry-notifier/internal/template"
)

func newNotifier(config cfg.Config) notifier.Notifier {
	renderer := template.NewRenderer(collectTemplates(config))

	return notifier.CreateNotifier(collectMessengers(config), renderer, config.Notify)
}

func collectTemplates(config cfg.Config) map[string]string {
	templates := map[string]string{}
	for _, tmpls := range config.Notify.On {
		for _, tmpl := range tmpls {
			templates[tmpl.MessageTemplateID] = tmpl.Message
		}
	}

	return templates
}

func collectMessengers(config cfg.Config) map[string][]messenger.Messenger {
	msgs := map[string][]messenger.Messenger{}

	for channelName, channel := range config.Channels {
		for _, tg := range channel.Telegram {
			if _, exists := msgs[channelName]; !exists {
				msgs[channelName] = []messenger.Messenger{}
			}

			msgs[channelName] = append(msgs[channelName], messenger.NewTelegram(*tg))
		}
	}

	return msgs
}
