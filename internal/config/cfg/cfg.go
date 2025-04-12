package cfg

import (
	"fmt"

	"github.com/artarts36/sentry-notifier/internal/messenger"
	"github.com/artarts36/sentry-notifier/internal/notifier"
	"github.com/artarts36/sentry-notifier/internal/security"
)

type Config struct {
	HTTP struct {
		Addr string `yaml:"address" json:"address"`
	} `yaml:"http" json:"http"`

	Control struct {
		Addr string `yaml:"address" json:"address"`
	} `yaml:"control" json:"control"`

	Log struct {
		Level string `yaml:"level" json:"level"`
	} `yaml:"log" json:"log"`

	Security security.Config `yaml:"security" json:"security"`

	Channels map[string]Channel `yaml:"channels" json:"channels"`

	Notify notifier.Config `yaml:"notify" json:"notify"`
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

	c.fillDefaults()

	return nil
}

func (c *Config) fillDefaults() {
	if c.HTTP.Addr == "" {
		c.HTTP.Addr = ":8080"
	}

	if c.Control.Addr == "" {
		c.Control.Addr = ":8081"
	}
}
