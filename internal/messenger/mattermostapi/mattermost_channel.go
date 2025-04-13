package mattermostapi

import (
	"context"
	"errors"
	"fmt"
	"github.com/artarts36/sentry-notifier/internal/messenger/errs"
	"net/http"
	"strings"

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
		find      func() (*model.Channel, *model.Response, error)
	}

	finders := []finder{
		{
			condition: m.cfg.Channel.ID == "" && m.cfg.Channel.Name == "",
			find: func() (*model.Channel, *model.Response, error) {
				return nil, nil, errors.New("must be set channel.id or channel.name")
			},
		},
		{
			condition: m.cfg.Channel.ID != "",
			find: func() (*model.Channel, *model.Response, error) {
				ch, resp, err := m.client.GetChannel(m.cfg.Channel.ID, "")
				if err != nil {
					return nil, resp, fmt.Errorf("find by id %q: %w", m.cfg.Channel.ID, err)
				}

				return ch, resp, nil
			},
		},
		{
			condition: m.cfg.Channel.TeamID != "",
			find: func() (*model.Channel, *model.Response, error) {
				ch, resp, err := m.client.GetChannelByName(m.cfg.Channel.Name, m.cfg.Channel.TeamID, "")
				if err != nil {
					return nil, resp, fmt.Errorf("find by name %q and team id %q: %w", m.cfg.Channel.Name, m.cfg.Channel.TeamID, err)
				}
				return ch, resp, nil
			},
		},
		{
			condition: m.cfg.Channel.TeamName != "",
			find: func() (*model.Channel, *model.Response, error) {
				ch, resp, err := m.client.GetChannelByNameForTeamName(m.cfg.Channel.Name, m.cfg.Channel.TeamName, "")
				if err != nil {
					return nil, resp, fmt.Errorf("find by name %q and team name %q: %w", m.cfg.Channel.Name, m.cfg.Channel.TeamName, err) //nolint:lll // not need
				}
				return ch, resp, nil
			},
		},
	}

	for _, f := range finders {
		if !f.condition {
			continue
		}

		channel, resp, err := f.find()
		if err == nil {
			return channel, nil
		}

		if resp != nil {
			return nil, m.mapChannelError(resp, err)
		}

		return nil, errs.NewNetworkError(err)
	}

	return nil, errors.New("channel misconfigured")
}

func (m *Mattermost) mapChannelError(resp *model.Response, err error) error {
	switch resp.StatusCode {
	case http.StatusInternalServerError:
		return errs.NewMessengerInternalError(err)
	case http.StatusNotFound:
		if strings.Contains(err.Error(), "Channel does not exist") {
			return errs.NewChatNotFoundError(err)
		}

		if strings.Contains(err.Error(), "Unable to find the existing team") {
			return errs.NewChatNotFoundErrorWithReason(err, "team_not_found")
		}

		return errs.NewUnexpectedError(err)
	case http.StatusUnauthorized:
		return errs.NewInvalidCredentialsError(err)
	default:
		return errs.NewUnexpectedError(err)
	}
}
