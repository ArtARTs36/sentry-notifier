package messenger

import (
	"errors"
)

const defaultTelegramHost = "api.telegram.org"

type TelegramConfig struct {
	Host     string `yaml:"host" json:"host"`
	BotToken string `yaml:"bot_token" json:"bot_token"`
	ChatID   string `yaml:"chat_id" json:"chat_id"`
	ThreadID string `yaml:"thread_id,omitempty" json:"thread_id"`
}

func (c *TelegramConfig) Validate() error {
	if c.Host == "" {
		c.Host = defaultTelegramHost
	}

	if c.BotToken == "" {
		return errors.New("bot_token must be set")
	}

	if c.BotToken == "" {
		return errors.New("bot_token must be set")
	}

	if c.ChatID == "" {
		return errors.New("chat_id must be set")
	}

	return nil
}
