package messenger

import (
	"context"
	"errors"
	"fmt"

	"github.com/mattermost/mattermost-server/v6/model"
)

func (m *Mattermost) retrieveChannel(ctx context.Context) error {
	if m.channel != nil {
		return nil
	}
	return m.loadChannel(ctx)
}

func (m *Mattermost) loadChannel(ctx context.Context) error {
	ch, err := m.findChannel(ctx)
	if err != nil {
		return fmt.Errorf("find: %w", err)
	}

	m.channel = ch

	return nil
}

func (m *Mattermost) findChannel(_ context.Context) (*model.Channel, error) {
	type finder struct {
		condition bool
		find      func() (*model.Channel, error)
	}

	finders := []finder{
		{
			condition: m.cfg.Channel.ID == "" && m.cfg.Channel.Name == "",
			find: func() (*model.Channel, error) {
				return nil, errors.New("must be set channel.id or channel.name")
			},
		},
		{
			condition: m.cfg.Channel.ID != "",
			find: func() (*model.Channel, error) {
				ch, _, err := m.client.GetChannel(m.cfg.Channel.ID, "")
				if err != nil {
					return nil, fmt.Errorf("find by id %q: %w", m.cfg.Channel.ID, err)
				}

				return ch, nil
			},
		},
		{
			condition: m.cfg.Channel.TeamID != "",
			find: func() (*model.Channel, error) {
				ch, _, err := m.client.GetChannelByName(m.cfg.Channel.Name, m.cfg.Channel.TeamID, "")
				if err != nil {
					return nil, fmt.Errorf("find by name %q and team id %q: %w", m.cfg.Channel.Name, m.cfg.Channel.TeamID, err)
				}
				return ch, nil
			},
		},
		{
			condition: m.cfg.Channel.TeamName != "",
			find: func() (*model.Channel, error) {
				ch, _, err := m.client.GetChannelByNameForTeamName(m.cfg.Channel.Name, m.cfg.Channel.TeamName, "")
				if err != nil {
					return nil, fmt.Errorf("find by name %q and team name %q: %w", m.cfg.Channel.Name, m.cfg.Channel.TeamName, err)
				}
				return ch, nil
			},
		},
	}

	for _, f := range finders {
		if f.condition {
			return f.find()
		}
	}

	return nil, errors.New("channel misconfigured")
}


