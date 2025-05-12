package middleware

import (
	"net/http"

	"github.com/slok/go-http-metrics/middleware"

	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	middlewarestd "github.com/slok/go-http-metrics/middleware/std"
)

func Metrics(target http.Handler) http.Handler {
	mdlw := middleware.New(middleware.Config{
		Service:  "webhook",
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	return middlewarestd.Handler("listen", mdlw, target)
}
