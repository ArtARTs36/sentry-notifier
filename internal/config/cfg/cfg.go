package cfg

import (
	"fmt"
	"github.com/artarts36/sentry-notifier/internal/messenger"
	"github.com/artarts36/sentry-notifier/internal/notifier"
	"github.com/artarts36/sentry-notifier/internal/security"
)

type Config struct {
	HTTP struct {
		Addr string `yaml:"address"`
	} `yaml:"http"`

	Control struct {
		Addr string `yaml:"address"`
	} `yaml:"control"`

	Log struct {
		Level string `yaml:"level"`
	} `yaml:"log"`

	Security security.Config `yaml:"security"`

	Channels map[string]Channel `yaml:"channels"`

	Notify notifier.Config `yaml:"notify"`
}

type Channel struct {
	Telegram   []*messenger.TelegramConfig   `yaml:"telegram,omitempty" json:"telegram,omitempty"`
	Mattermost []*messenger.MattermostConfig `yaml:"mattermost,omitempty" json:"mattermost,omitempty"`
}

func (c *Channel) Validate() error {
	for i, mm := range c.Mattermost {
		if err := mm.Validate(); err != nil {
			return fmt.Errorf("mattermost[%d]: %w", i, err)
		}
	}

	return nil
}

func (c *Config) Validate() error {
	for channelName, channel := range c.Channels {
		if err := channel.Validate(); err != nil {
			return fmt.Errorf("channel %q: %w", channelName, err)
		}
	}

	return nil
}
