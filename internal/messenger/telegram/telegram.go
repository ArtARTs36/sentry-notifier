package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/artarts36/sentry-notifier/internal/messenger/contracts"
	"github.com/artarts36/sentry-notifier/internal/messenger/errs"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

type Telegram struct {
	cfg Config
}

type telegramRequest struct {
	ChatID          string `json:"chat_id"`
	Text            string `json:"text"`
	MessageThreadID string `json:"message_thread_id"`
	ParseMode       string `json:"parse_mode"`
}

type telegramErrorResponse struct {
	OK          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

func NewTelegram(cfg Config) *Telegram {
	return &Telegram{
		cfg: cfg,
	}
}

func (t *Telegram) Name() string {
	return "telegram"
}

func (t *Telegram) Ping(ctx context.Context) error {
	return t.getChat(ctx)
}

func (t *Telegram) Send(ctx context.Context, message contracts.Message) error {
	msgBody := strings.ReplaceAll(message.Body, "-", "\\-")
	msgBody = strings.ReplaceAll(msgBody, ".", "\\.")

	payload, err := json.Marshal(&telegramRequest{
		ChatID:          t.cfg.ChatID.Value,
		MessageThreadID: t.cfg.ThreadID.Value,
		Text:            msgBody,
		ParseMode:       "MarkdownV2",
	})
	if err != nil {
		return fmt.Errorf("failed to marshal telegram request: %w", err)
	}
	body := bytes.NewBuffer(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.buildURL(), body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}

	defer func() {
		if bodyCloseErr := resp.Body.Close(); bodyCloseErr != nil {
			slog.
				With(slog.Any("err", bodyCloseErr)).
				ErrorContext(ctx, "[telegram] failed to close response body")
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("telegram: failed to read response body: %w", err)
	}

	respLogLevel := slog.LevelDebug
	if resp.StatusCode != http.StatusOK {
		respLogLevel = slog.LevelWarn
	}

	slog.
		With(slog.Int("status_code", resp.StatusCode)).
		With(slog.Any("body", string(respBody))).
		Log(ctx, respLogLevel, "[telegram] received response")

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram: failed to send message: got %d status code", resp.StatusCode)
	}

	return nil
}

func (t *Telegram) buildURL() string {
	return fmt.Sprintf(
		"https://%s/bot%s/sendMessage",
		t.cfg.Host.Value,
		t.cfg.BotToken.Value,
	)
}

func (t *Telegram) getChat(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.buildGetChatURL(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errs.NewNetworkError(err)
	}

	defer func() {
		if bodyCloseErr := resp.Body.Close(); bodyCloseErr != nil {
			slog.
				With(slog.Any("err", bodyCloseErr)).
				ErrorContext(ctx, "[telegram] failed to close response body")
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("telegram: failed to read response body: %w", err)
	}

	respLogLevel := slog.LevelDebug
	if resp.StatusCode != http.StatusOK {
		respLogLevel = slog.LevelWarn
	}

	slog.
		With(slog.Int("status_code", resp.StatusCode)).
		With(slog.Any("body", string(respBody))).
		Log(ctx, respLogLevel, "[telegram] received get chat response")

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	return t.mapChatError(resp, respBody)
}

func (t *Telegram) mapChatError(resp *http.Response, respBody []byte) error {
	var errResp telegramErrorResponse

	err := json.Unmarshal(respBody, &errResp)
	if err != nil {
		return errs.NewUnexpectedError(fmt.Errorf("unmarshal response json: %w", err))
	}

	newError := func() error {
		return fmt.Errorf("telegram returns: %s", errResp.Description)
	}

	switch resp.StatusCode {
	case http.StatusBadRequest:
		if strings.Contains(errResp.Description, "chat not found") {
			return errs.NewChatNotFoundError(newError())
		}
	case http.StatusUnauthorized:
		return errs.NewInvalidCredentialsError(newError())
	case http.StatusInternalServerError:
		return errs.NewMessengerInternalError(newError())
	}

	return errs.NewUnexpectedError(newError())
}

func (t *Telegram) buildGetChatURL() string {
	return fmt.Sprintf(
		"https://%s/bot%s/getChat?chat_id=%s",
		t.cfg.Host.Value,
		t.cfg.BotToken.Value,
		t.cfg.ChatID.Value,
	)
}
