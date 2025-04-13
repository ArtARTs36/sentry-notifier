package messenger

import (
	"context"
	"errors"
	"fmt"
	"github.com/mattermost/mattermost-server/v6/model"
)

type Mattermost struct {
	cfg     MattermostAPIConfig
	client  *model.Client4
	channel *model.Channel
}

type MattermostConfig struct {
	API *MattermostAPIConfig `yaml:"api" json:"api"`
}

type MattermostAPIConfig struct {
	Token  string `yaml:"token" json:"token"`
	Server string `yaml:"server" json:"server"`

	Channel struct {
		ID string `yaml:"id" json:"channel_id"` // one of

		Name string `yaml:"name" json:"channel_name"`

		TeamName string `yaml:"team_name" json:"team"`
		TeamID   string `yaml:"team_id" json:"team_id"`
	} `yaml:"channel" json:"channel"`
}

func (c *MattermostConfig) Validate() error {
	if c.API == nil {
		return errors.New("must be set api configuration")
	}

	if c.API.Channel.ID != "" {
		return nil
	}

	if c.API.Channel.Name == "" {
		return errors.New("must be set channel.id or channel.name")
	}

	if c.API.Channel.TeamID == "" && c.API.Channel.TeamName == "" {
		return errors.New("when use channel name, must be set team_name or team_id")
	}

	return nil
}

func NewMattermost(cfg MattermostConfig) *Mattermost {
	client := model.NewAPIv4Client(cfg.API.Server)
	client.SetToken(cfg.API.Token)

	m := &Mattermost{
		cfg:    *cfg.API,
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
	return "mattermost"
}

func (m *Mattermost) Send(ctx context.Context, message Message) error {
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
