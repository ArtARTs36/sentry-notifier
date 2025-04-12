package health

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

func Handler(checker func(ctx context.Context) *Check) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		result := checker(r.Context())

		body, err := json.Marshal(result)
		if err != nil {
			slog.ErrorContext(r.Context(), "failed to marshal health check response", slog.Any("err", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if result.Status {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		_, err = w.Write(body)
		if err != nil {
			slog.ErrorContext(r.Context(), "failed to write response", slog.Any("err", err))
			return
		}
	}
}
