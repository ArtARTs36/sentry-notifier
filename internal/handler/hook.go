package handler

import (
	"github.com/artarts36/sentry-notifier/internal/notifier"
	"github.com/artarts36/sentry-notifier/internal/sentry"
	"io"
	"log/slog"
	"net/http"
)

type HookHandler struct {
	notifier notifier.Notifier
}

func NewHookHandler(notifier notifier.Notifier) *HookHandler {
	return &HookHandler{
		notifier: notifier,
	}
}

func (h *HookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	eventName, err := sentry.WrapHookResource(r.Header.Get(sentry.HeaderHookResource))
	if err != nil {
		slog.
			With(slog.String("err", err.Error())).
			ErrorContext(r.Context(), "failed to wrap hook resource")

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rawPayload, err := io.ReadAll(r.Body)
	if err != nil {
		slog.
			With(slog.String("err", err.Error())).
			ErrorContext(r.Context(), "[hook] failed to read payload")

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	payload, err := sentry.ParsePayload(eventName, rawPayload)
	if err != nil {
		slog.
			With(slog.String("err", err.Error())).
			ErrorContext(r.Context(), "[hook] failed to parse payload")

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	err = h.notifier.Notify(r.Context(), payload)
	if err != nil {
		slog.
			With(slog.String("err", err.Error())).
			ErrorContext(r.Context(), "[hook] failed to notify")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
