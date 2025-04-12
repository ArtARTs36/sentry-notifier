package messenger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

const defaultTelegramHost = "api.telegram.org"

type Telegram struct {
	cfg TelegramConfig
}

type TelegramConfig struct {
	Host     string `yaml:"host" json:"host"`
	BotToken string `yaml:"bot_token" json:"bot_token"`
	ChatID   string `yaml:"chat_id" json:"chat_id"`
	ThreadID string `yaml:"thread_id,omitempty" json:"thread_id"`
}

type telegramRequest struct {
	ChatID          string `json:"chat_id"`
	Text            string `json:"text"`
	MessageThreadID string `json:"message_thread_id"`
	ParseMode       string `json:"parse_mode"`
}

func NewTelegram(cfg TelegramConfig) *Telegram {
	if cfg.Host == "" {
		cfg.Host = defaultTelegramHost
	}

	return &Telegram{
		cfg: cfg,
	}
}

func (t *Telegram) Name() string {
	return "telegram"
}

func (t *Telegram) Ping(_ context.Context) error {
	return nil
}

func (t *Telegram) Send(ctx context.Context, message Message) error {
	msgBody := strings.ReplaceAll(message.Body, "-", "\\-")
	msgBody = strings.ReplaceAll(msgBody, ".", "\\.")

	payload, err := json.Marshal(&telegramRequest{
		ChatID:          t.cfg.ChatID,
		MessageThreadID: t.cfg.ThreadID,
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
		return fmt.Errorf("failed to send message: %w", err)
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
	return fmt.Sprintf("https://%s/bot%s/sendMessage", t.cfg.Host, t.cfg.BotToken)
}
