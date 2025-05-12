package cfg

import (
	"fmt"
	"github.com/artarts36/sentry-notifier/internal/messenger/mattermostwh"
	"github.com/artarts36/sentry-notifier/internal/port/http/middleware"
	"github.com/artarts36/specw"

	"github.com/artarts36/sentry-notifier/internal/messenger/mattermostapi"
	"github.com/artarts36/sentry-notifier/internal/messenger/telegram"
	"github.com/artarts36/sentry-notifier/internal/notifier"
)

type Config struct {
	HTTP struct {
		Addr      string                     `yaml:"address" json:"address"`
		RateLimit middleware.RateLimitConfig `yaml:"rate_limit" json:"rate_limit"`
	} `yaml:"http" json:"http"`

	Control struct {
		Addr string `yaml:"address" json:"address"`
	} `yaml:"control" json:"control"`

	Log struct {
		Level specw.SlogLevel `yaml:"level" json:"level"`
	} `yaml:"log" json:"log"`

	Security middleware.AuthorizeConfig `yaml:"security" json:"security"`

	Channels map[string]Channel `yaml:"channels" json:"channels"`

	Notify notifier.Config `yaml:"notify" json:"notify"`
}

type Channel struct {
	Telegram          []*telegram.Config      `yaml:"telegram,omitempty" json:"telegram,omitempty"`
	MattermostAPI     []*mattermostapi.Config `yaml:"mattermost_api,omitempty" json:"mattermost_api,omitempty"`
	MattermostWebhook []*mattermostwh.Config  `yaml:"mattermost_webhook,omitempty" json:"mattermost_webhook,omitempty"`
}

func (c *Channel) Validate() error {
	for i, mm := range c.MattermostAPI {
		if err := mm.Validate(); err != nil {
			return fmt.Errorf("mattermost_api[%d]: %w", i, err)
		}
	}

	for i, tg := range c.Telegram {
		if err := tg.Validate(); err != nil {
			return fmt.Errorf("telegram[%d]: %w", i, err)
		}
	}

	return nil
}

func (c *Config) Validate() error {
	if err := c.HTTP.RateLimit.Validate(); err != nil {
		return fmt.Errorf("rate_limit: %w", err)
	}

	for channelName, channel := range c.Channels {
		if err := channel.Validate(); err != nil {
			return fmt.Errorf("channel[%q]: %w", channelName, err)
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
