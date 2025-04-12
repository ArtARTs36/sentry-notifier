package app

import (
	"context"
	goMetrics "github.com/artarts36/go-metrics"
	"github.com/artarts36/sentry-notifier/internal/handler"
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

	metricsRegistry goMetrics.Registry
}

func New(config cfg.Config, metricsRegistry goMetrics.Registry) *Server {
	notif := newNotifier(config)

	s := &Server{
		config:          config,
		notifier:        notif,
		metricsRegistry: metricsRegistry,
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
