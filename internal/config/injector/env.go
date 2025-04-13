package injector

import (
	"errors"
	"os"

	"github.com/artarts36/sentry-notifier/internal/config/cfg"
)

type Env struct {
}

func NewEnv() *Env {
	return &Env{}
}

func (e *Env) Inject(config cfg.Config) (cfg.Config, error) { //nolint:gocognit // complex
	for _, channel := range config.Channels {
		for _, telegramConfig := range channel.Telegram {
			if channel.Telegram != nil {
				var err error
				telegramConfig.BotToken, err = e.transform(telegramConfig.BotToken)
				if err != nil {
					return cfg.Config{}, err
				}
				telegramConfig.ChatID, err = e.transform(telegramConfig.ChatID)
				if err != nil {
					return cfg.Config{}, err
				}
				telegramConfig.ThreadID, err = e.transform(telegramConfig.ThreadID)
				if err != nil {
					return cfg.Config{}, err
				}
			}
		}

		for _, mm := range channel.MattermostAPI {
			if mm != nil {
				var err error
				mm.Server, err = e.transform(mm.Server)
				if err != nil {
					return cfg.Config{}, err
				}
				mm.Token, err = e.transform(mm.Token)
				if err != nil {
					return cfg.Config{}, err
				}
			}
		}
	}

	var err error
	config.Security.ClientToken, err = e.transform(config.Security.ClientToken)
	if err != nil {
		return cfg.Config{}, err
	}

	return config, nil
}

func (e *Env) transform(stringWithVar string) (string, error) {
	if stringWithVar == "" {
		return stringWithVar, nil
	}

	if stringWithVar[0] != '$' {
		return stringWithVar, nil
	}

	varName := stringWithVar[1:]
	varValue := os.Getenv(varName)
	if varValue == "" {
		return "", errors.New(varName + " is not set")
	}

	return varValue, nil
}
