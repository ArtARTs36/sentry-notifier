package app

import (
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	middlewarestd "github.com/slok/go-http-metrics/middleware/std"
	"log/slog"
	"net/http"

	"github.com/artarts36/sentry-notifier/internal/config/cfg"
	appmw "github.com/artarts36/sentry-notifier/internal/port/http/middleware"
	sloghttp "github.com/samber/slog-http"
)

func (s *Server) buildHandler(target http.Handler, config cfg.Config) http.Handler {
	return sloghttp.New(slog.Default())(
		wrapMetricHandler(
			appmw.AuthorizeRequest(
				target,
				config.Security,
			),
		),
	)
}

func wrapMetricHandler(target http.Handler) http.Handler {
	mdlw := middleware.New(middleware.Config{
		Service:  "webhook",
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	return middlewarestd.Handler("listen", mdlw, target)
}
