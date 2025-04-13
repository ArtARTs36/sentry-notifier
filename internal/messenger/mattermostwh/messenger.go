package mattermostwh

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/artarts36/sentry-notifier/internal/messenger/contracts"
	"io"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/artarts36/sentry-notifier/internal/messenger/errs"
)

type Messenger struct {
	cfg Config

	host string
	url  string
}

type Config struct {
	URL urlWrapper `yaml:"url" json:"url"`
}

func NewMessenger(cfg Config) *Messenger {
	return &Messenger{
		cfg:  cfg,
		host: fmt.Sprintf("%s//%s", cfg.URL.Scheme, cfg.URL.Host),
		url:  cfg.URL.String(),
	}
}

func (m *Messenger) Name() string {
	return "mattermost_webhook"
}

func (m *Messenger) Ping(ctx context.Context) error {
	const defaultTimeout = 30 * time.Second

	d := net.Dialer{Timeout: defaultTimeout}

	conn, err := d.DialContext(ctx, "tcp", m.host)
	if err != nil {
		return errs.NewNetworkError(err)
	}

	_ = conn.Close()

	return nil
}

type mattermostWebhookRequest struct {
	Text string `json:"text"`
}

func (m *Messenger) Send(ctx context.Context, message contracts.Message) error {
	body, err := m.marshallMessage(message)
	if err != nil {
		return fmt.Errorf("encode message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, m.url, body)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}

	defer func() {
		if bodyCloseErr := resp.Body.Close(); bodyCloseErr != nil {
			slog.
				With(slog.Any("err", bodyCloseErr)).
				ErrorContext(ctx, "[mattermost-webhook] failed to close response body")
		}
	}()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return errs.NewInvalidCredentialsError(errors.New("401 Unauthorized"))
	}

	return errs.NewUnexpectedError(fmt.Errorf("mattermost returns %d status", resp.StatusCode))
}

func (m *Messenger) marshallMessage(message contracts.Message) (io.Reader, error) {
	req := &mattermostWebhookRequest{
		Text: message.Body,
	}

	content, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(content), nil
}
