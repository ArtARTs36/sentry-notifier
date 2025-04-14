package telegram

import (
	"errors"
	"github.com/artarts36/specw"
)

const defaultTelegramHost = "api.telegram.org"

type Config struct {
	Host     specw.Env[string] `yaml:"host" json:"host"`
	BotToken specw.Env[string] `yaml:"bot_token" json:"bot_token"`
	ChatID   specw.Env[string] `yaml:"chat_id" json:"chat_id"`
	ThreadID specw.Env[string] `yaml:"thread_id,omitempty" json:"thread_id"`
}

func (c *Config) Validate() error {
	if c.Host.Value == "" {
		c.Host.Value = defaultTelegramHost
	}

	if c.BotToken.Value == "" {
		return errors.New("bot_token must be set")
	}

	if c.BotToken.Value == "" {
		return errors.New("bot_token must be set")
	}

	if c.ChatID.Value == "" {
		return errors.New("chat_id must be set")
	}

	return nil
}
