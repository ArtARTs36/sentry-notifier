package app

import (
	"context"
	goMetrics "github.com/artarts36/go-metrics"
	"github.com/artarts36/sentry-notifier/internal/handler"
	"github.com/artarts36/sentry-notifier/internal/health"
	"github.com/artarts36/sentry-notifier/internal/messenger"
	"github.com/artarts36/sentry-notifier/internal/metrics"
	"github.com/artarts36/sentry-notifier/internal/notifier"
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

	messengers map[string][]messenger.Messenger
}

func New(config cfg.Config, metricsRegistry goMetrics.Registry) *Server {
	metr := metrics.NewGroup(metricsRegistry)

	notif, messengers := newNotifier(config, metr)

	s := &Server{
		config:     config,
		notifier:   notif,
		messengers: messengers,
	}

	s.handler = s.buildHandler(handler.NewHookHandler(notif), config)

	return s
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

	s.notifier.Close()

	return err
}

func (s *Server) Health(ctx context.Context) *health.Check {
	result := &health.Check{
		Status:   true,
		Channels: map[string]map[string][]health.CheckChannel{},
	}

	reason := func(err error) string {
		if err != nil {
			return err.Error()
		}
		return ""
	}

	for channelName, ms := range s.messengers {
		result.Channels[channelName] = map[string][]health.CheckChannel{}

		for _, m := range ms {
			pingErr := m.Ping(ctx)
			if pingErr != nil {
				result.Status = false
			}

			result.Channels[channelName][m.Name()] = append(result.Channels[channelName][m.Name()], health.CheckChannel{
				Status: pingErr == nil,
				Reason: reason(pingErr),
			})
		}
	}

	return result
}
