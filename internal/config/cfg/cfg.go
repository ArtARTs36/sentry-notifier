package cfg

import (
	"github.com/artarts36/sentry-notifier/internal/messenger"
	"github.com/artarts36/sentry-notifier/internal/notifier"
	"github.com/artarts36/sentry-notifier/internal/security"
)

type Config struct {
	HTTP struct {
		Addr string `yaml:"address"`
	} `yaml:"http"`

	Metrics struct {
		Addr string `yaml:"address"`
	} `yaml:"metrics"`

	Log struct {
		Level string `yaml:"level"`
	} `yaml:"log"`

	Security security.Config `yaml:"security"`

	Channels map[string]Channel `yaml:"channels"`

	Notify notifier.Config `yaml:"notify"`
}

type Channel struct {
	Telegram []*messenger.TelegramConfig `yaml:"telegram,omitempty"`
}
