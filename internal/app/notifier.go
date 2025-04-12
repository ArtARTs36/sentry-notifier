package app

import (
	"github.com/artarts36/sentry-notifier/internal/config/cfg"
	"github.com/artarts36/sentry-notifier/internal/messenger"
	"github.com/artarts36/sentry-notifier/internal/metrics"
	"github.com/artarts36/sentry-notifier/internal/notifier"
	"github.com/artarts36/sentry-notifier/internal/template"
)

func newNotifier(config cfg.Config, metr *metrics.Group) (notifier.Notifier, map[string][]messenger.Messenger) {
	renderer := template.NewRenderer(collectTemplates(config))
	messengers := collectMessengers(config)

	return notifier.CreateNotifier(messengers, renderer, config.Notify, metr), messengers
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
		if _, exists := msgs[channelName]; !exists {
			msgs[channelName] = []messenger.Messenger{}
		}

		for _, tg := range channel.Telegram {
			msgs[channelName] = append(msgs[channelName], messenger.NewTelegram(*tg))
		}

		for _, mm := range channel.Mattermost {
			msgs[channelName] = append(msgs[channelName], messenger.NewMattermost(*mm))
		}
	}

	return msgs
}
