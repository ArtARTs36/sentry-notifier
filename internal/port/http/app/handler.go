package app

import (
	"log/slog"
	"net/http"

	"github.com/artarts36/sentry-notifier/internal/config/cfg"
	appmw "github.com/artarts36/sentry-notifier/internal/port/http/middleware"
	sloghttp "github.com/samber/slog-http"
)

func (s *Server) buildHandler(target http.Handler, config cfg.Config) http.Handler {
	return appmw.RateLimit(
		appmw.Pattern("/listen",
			sloghttp.New(slog.Default())(
				appmw.Metrics(
					appmw.AuthorizeRequest(
						target,
						config.Security,
					),
				),
			),
		),
		config.HTTP.RateLimit,
	)
}
