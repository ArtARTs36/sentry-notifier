package notifier

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/tyler-sommer/stick"

	"github.com/artarts36/sentry-notifier/internal/messenger"
	"github.com/artarts36/sentry-notifier/internal/metrics"
	"github.com/artarts36/sentry-notifier/internal/sentry"
	"github.com/artarts36/sentry-notifier/internal/template"
)

type ImmediatelyNotifier struct {
	messengers map[string][]messenger.Messenger
	renderer   *template.Renderer
	cfg        Config
	metrics    *metrics.Messages
}

type Template struct {
	Message string `yaml:"message" json:"message"`
	To      string `yaml:"to" json:"to"`

	MessageTemplateID string `yaml:"-" json:"-"`
}

type messengersMessage struct {
	messengers  []messenger.Messenger
	message     string
	channelName string
}

func NewImmediatelyNotifier(
	messengers map[string][]messenger.Messenger,
	renderer *template.Renderer,
	cfg Config,
	msgMetrics *metrics.Messages,
) *ImmediatelyNotifier {
	return &ImmediatelyNotifier{
		messengers: messengers,
		renderer:   renderer,
		cfg:        cfg,
		metrics:    msgMetrics,
	}
}

func (n *ImmediatelyNotifier) Notify(ctx context.Context, payload sentry.Payload) error {
	messengersMessages, allMessagesCount, err := n.prepareMessages(ctx, payload)
	if err != nil {
		if errors.Is(err, errTemplatesNotSelected) {
			slog.DebugContext(ctx, "[notifier] no found templates")

			return nil
		}

		return fmt.Errorf("prepare messages: %w", err)
	}

	slog.InfoContext(ctx, fmt.Sprintf("[notifier] sending %d messages", allMessagesCount))
	wg := &sync.WaitGroup{}

	for _, tmpl := range messengersMessages {
		for _, mess := range tmpl.messengers {
			wg.Add(1)

			go func() {
				defer wg.Done()

				slog.
					InfoContext(ctx, "[notifier] sending message", slog.String("messenger", mess.Name()))

				n.metrics.IncSending(string(payload.GetHookResource()), mess.Name(), tmpl.channelName)

				serr := mess.Send(ctx, messenger.Message{
					Body: tmpl.message,
				})
				if serr != nil {
					slog.
						With(slog.String("messenger", mess.Name())).
						With(slog.Any("err", serr)).
						WarnContext(ctx, "failed to send message")
				}

				n.metrics.IncSent(string(payload.GetHookResource()), mess.Name(), tmpl.channelName, serr == nil)
			}()
		}
	}

	wg.Wait()

	return nil
}

var errTemplatesNotSelected = errors.New("no templates selected")

func (n *ImmediatelyNotifier) prepareMessages(
	ctx context.Context,
	payload sentry.Payload,
) ([]messengersMessage, int, error) {
	messengersMessages := make([]messengersMessage, 0)

	templates := n.selectTemplates(payload)
	if len(templates) == 0 {
		return nil, 0, errTemplatesNotSelected
	}

	slog.DebugContext(ctx, fmt.Sprintf("[notifier] selected %d templates", len(templates)))

	allMessagesCount := 0

	for _, tmpl := range templates {
		msgs, lErr := n.selectMessengers(tmpl.To)
		if lErr != nil {
			return nil, 0, fmt.Errorf("notifier: failed to select messengersMessages: %w", lErr)
		}

		message, lErr := n.renderer.Render(tmpl.MessageTemplateID, map[string]stick.Value{
			"hook": payload.GetData(),
		})
		if lErr != nil {
			return nil, 0, fmt.Errorf("notifier: failed to render message: %w", lErr)
		}

		messengersMessages = append(messengersMessages, messengersMessage{
			messengers:  msgs,
			message:     string(message),
			channelName: tmpl.To,
		})

		allMessagesCount += len(msgs)
	}

	return messengersMessages, allMessagesCount, nil
}

func (n *ImmediatelyNotifier) selectMessengers(channelName string) ([]messenger.Messenger, error) {
	msgs, exists := n.messengers[channelName]
	if !exists {
		return nil, fmt.Errorf("notifier: no found messengers for channel %s", channelName)
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
