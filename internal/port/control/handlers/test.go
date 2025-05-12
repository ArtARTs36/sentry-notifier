package handlers

import (
	"encoding/json"
	"github.com/artarts36/sentry-notifier/internal/sentry"
	"log/slog"
	"net/http"

	"github.com/artarts36/sentry-notifier/internal/notifier"
)

type TestHandler struct {
	notifier notifier.Notifier
}

func NewTestHandler(notifier notifier.Notifier) *TestHandler {
	return &TestHandler{
		notifier: notifier,
	}
}

type sendTestNotificationResponse struct {
	Message string `json:"message"`
}

func (h *TestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var payload sentry.Payload = sentry.ExampleIssuePayload()
	if eventType := r.URL.Query().Get("event"); eventType == sentry.HookResourceEventAlert {
		payload = sentry.ExampleEventAlert()
	}

	err := h.notifier.Notify(r.Context(), payload)
	if err != nil {
		slog.
			With(slog.Any("err", err)).
			ErrorContext(r.Context(), "[hook] failed to notify")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := sendTestNotificationResponse{}

	switch h.notifier.(type) {
	case *notifier.ImmediatelyNotifier:
		resp.Message = "Notifications was sent immediately"
	case *notifier.AsyncNotifier:
		resp.Message = "Notifications was scheduled"
	case *notifier.NoopNotifier:
		resp.Message = "Used null notify strategy. Notifications didnt sent"
	}

	body, err := json.Marshal(resp)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to marshal send test notification response", slog.Any("err", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(body)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to write response", slog.Any("err", err))
		return
	}
}
