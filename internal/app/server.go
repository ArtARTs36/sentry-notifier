package app

import (
	"context"
	"errors"
	"fmt"
	goMetrics "github.com/artarts36/go-metrics"
	"github.com/artarts36/sentry-notifier/internal/handler"
	"github.com/artarts36/sentry-notifier/internal/health"
	"github.com/artarts36/sentry-notifier/internal/messenger/contracts"
	messengererrors "github.com/artarts36/sentry-notifier/internal/messenger/errs"
	"github.com/artarts36/sentry-notifier/internal/metrics"
	"github.com/artarts36/sentry-notifier/internal/notifier"
	"log/slog"
	"net/http"
	"time"

	"github.com/artarts36/sentry-notifier/internal/config/cfg"
)

const httpReadTimeout = 10 * time.Second

type Server struct {
	config cfg.Config

	handler http.Handler

	server *http.Server

	notifier notifier.Notifier

	messengers map[string][]contracts.Messenger
}

func New(config cfg.Config, metricsRegistry goMetrics.Registry) (*Server, notifier.Notifier) {
	metr := metrics.NewGroup(metricsRegistry)

	notif, messengers := newNotifier(config, metr)

	s := &Server{
		config:     config,
		notifier:   notif,
		messengers: messengers,
	}

	s.handler = s.buildHandler(handler.NewHookHandler(notif), config)

	return s, notif
}

func (s *Server) Run() error {
	mux := http.NewServeMux()
	mux.Handle("/listen", s.handler)

	hServer := &http.Server{
		Addr:        s.config.HTTP.Addr,
		Handler:     mux,
		ReadTimeout: httpReadTimeout,
	}

	s.server = hServer

	return hServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	err := s.server.Shutdown(ctx)

	return err
}

func (s *Server) Health(ctx context.Context) *health.CheckResponse {
	result := &health.CheckResponse{
		Status: true,
		Checks: make([]health.Check, 0),
	}

	reason := func(err error) string {
		if err != nil {
			var perr messengererrors.Error
			if errors.As(err, &perr) {
				return perr.Reason()
			}

			return "unexpected"
		}

		return ""
	}

	for channelName, ms := range s.messengers {
		for i, m := range ms {
			pingErr := m.Ping(ctx)
			pingErrReason := reason(pingErr)
			if pingErr != nil {
				result.Status = false
				slog.ErrorContext(
					ctx,
					"[health] failed to ping messenger",
					slog.String("messenger", m.Name()),
					slog.String("channel", channelName),
					slog.Any("err", pingErr),
					slog.String("err_reason", pingErrReason),
				)
			}

			result.Checks = append(result.Checks, health.Check{
				ID:     fmt.Sprintf("%s-%s-%d", channelName, m.Name(), i),
				Status: pingErr == nil,
				Reason: pingErrReason,
			})
		}
	}

	return result
}
