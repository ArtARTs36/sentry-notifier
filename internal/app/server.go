package app

import (
	"context"
	"github.com/artarts36/sentry-notifier/internal/handler"
	"github.com/artarts36/sentry-notifier/internal/notifier"
	"github.com/artarts36/sentry-notifier/internal/security"
	sloghttp "github.com/samber/slog-http"
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
}

func New(config cfg.Config) *Server {
	notif := newNotifier(config)

	return &Server{
		config: config,

		handler: sloghttp.New(slog.Default())(security.AuthorizeRequest(
			handler.NewHookHandler(notif),
			config.Security,
		)),

		notifier: notif,
	}
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
