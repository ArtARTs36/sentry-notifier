package notifier

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/tyler-sommer/stick"

	"github.com/artarts36/sentry-notifier/internal/messenger"
	"github.com/artarts36/sentry-notifier/internal/sentry"
	"github.com/artarts36/sentry-notifier/internal/template"
)

type ImmediatelyNotifier struct {
	messengers map[string][]messenger.Messenger
	renderer   *template.Renderer
	cfg        Config
}

type Template struct {
	Message string `yaml:"message"`
	To      string `yaml:"to"`

	MessageTemplateID string `yaml:"-"`
}

type messengersMessage struct {
	messengers []messenger.Messenger
	message    string
}

func NewImmediatelyNotifier(
	messengers map[string][]messenger.Messenger,
	renderer *template.Renderer,
	cfg Config,
) *ImmediatelyNotifier {
	return &ImmediatelyNotifier{
		messengers: messengers,
		renderer:   renderer,
		cfg:        cfg,
	}
}

func (n *ImmediatelyNotifier) Notify(ctx context.Context, payload sentry.Payload) error {
	templates := n.selectTemplates(payload)
	if len(templates) == 0 {
		slog.DebugContext(ctx, "[notifier] no templates selected")

		return nil
	}

	slog.DebugContext(ctx, fmt.Sprintf("[notifier] selected %d templates", len(templates)))

	messengersMessages := make([]messengersMessage, 0)
	wg := &sync.WaitGroup{}

	allMessagesCount := 0

	for _, tmpl := range templates {
		msgs, lErr := n.selectMessengers(tmpl.To)
		if lErr != nil {
			return fmt.Errorf("notifier: failed to select messengersMessages: %w", lErr)
		}

		message, lErr := n.renderer.Render(tmpl.MessageTemplateID, map[string]stick.Value{
			"hook": payload.GetData(),
		})
		if lErr != nil {
			return fmt.Errorf("notifier: failed to render message: %w", lErr)
		}

		messengersMessages = append(messengersMessages, messengersMessage{
			messengers: msgs,
			message:    string(message),
		})

		allMessagesCount += len(msgs)
	}

	slog.InfoContext(ctx, fmt.Sprintf("[notifier] sending %d messages", allMessagesCount))

	for _, tmpl := range messengersMessages {
		for _, mess := range tmpl.messengers {
			wg.Add(1)

			go func() {
				defer wg.Done()

				slog.
					InfoContext(ctx, "[notifier] sending message via %s", slog.String("messenger", mess.Name()))

				err := mess.Send(ctx, messenger.Message{
					Body: tmpl.message,
				})
				if err != nil {
					slog.
						With(slog.String("messenger", mess.Name())).
						With(slog.Any("err", err)).
						WarnContext(ctx, "failed to send message")
				}
			}()
		}
	}

	wg.Wait()

	return nil
}

func (n *ImmediatelyNotifier) selectMessengers(channelName string) ([]messenger.Messenger, error) {
	msgs, exists := n.messengers[channelName]
	if !exists {
		return nil, fmt.Errorf("notifier: failed to find messengers for channel %s", channelName)
	}

	return msgs, nil
}

func (n *ImmediatelyNotifier) selectTemplates(payload sentry.Payload) []Template {
	res := payload.GetHookResource()

	msgs, exists := n.cfg.On[res]
	if !exists {
		return []Template{}
	}

	return msgs
}

func (*ImmediatelyNotifier) Close() {
}
