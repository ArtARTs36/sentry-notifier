package mattermostapi

import (
	"context"
	"errors"
	"fmt"
	"github.com/artarts36/sentry-notifier/internal/messenger/contracts"
	"github.com/artarts36/specw"
	"github.com/mattermost/mattermost-server/v6/model"
)

type Mattermost struct {
	cfg     Config
	client  *model.Client4
	channel *model.Channel
}

type Config struct {
	Token  specw.Env[string] `yaml:"token" json:"token"`
	Server specw.Env[string] `yaml:"server" json:"server"`

	Channel struct {
		ID specw.Env[string] `yaml:"id" json:"channel_id"` // one of

		Name specw.Env[string] `yaml:"name" json:"channel_name"`

		TeamName specw.Env[string] `yaml:"team_name" json:"team"`
		TeamID   specw.Env[string] `yaml:"team_id" json:"team_id"`
	} `yaml:"channel" json:"channel"`
}

func (c *Config) Validate() error {
	if c.Channel.ID.Value != "" {
		return nil
	}

	if c.Channel.Name.Value == "" {
		return errors.New("must be set channel.id or channel.name")
	}

	if c.Channel.TeamID.Value == "" && c.Channel.TeamName.Value == "" {
		return errors.New("when use channel name, must be set team_name or team_id")
	}

	return nil
}

func NewMessenger(cfg Config) *Mattermost {
	client := model.NewAPIv4Client(cfg.Server.Value)
	client.SetToken(cfg.Token.Value)

	m := &Mattermost{
		cfg:    cfg,
		client: client,
	}

	return m
}

func (m *Mattermost) Ping(ctx context.Context) error {
	err := m.loadChannel(ctx)
	if err != nil {
		return fmt.Errorf("load channel: %w", err)
	}
	return nil
}

func (m *Mattermost) Name() string {
	return "mattermost_api"
}

func (m *Mattermost) Send(ctx context.Context, message contracts.Message) error {
	err := m.retrieveChannel(ctx)
	if err != nil {
		return fmt.Errorf("retrieve channel: %w", err)
	}

	_, _, err = m.client.CreatePost(&model.Post{
		ChannelId: m.channel.Id,
		Message:   message.Body,
	})
	if err != nil {
		return fmt.Errorf("create post: %w", err)
	}

	return nil
}
