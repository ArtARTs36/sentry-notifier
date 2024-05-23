package security

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log/slog"
	"net/http"
)

type Config struct {
	ClientToken string `yaml:"client_secret"`
}

func AuthorizeRequest(next http.Handler, cfg Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		next := next

		expectedDigest := request.Header.Get("Sentry-Hook-Signature")
		if expectedDigest == "" {
			slog.WarnContext(request.Context(), "[security] unauthorized request: empty Sentry-Hook-Signature")

			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if resource := request.Header.Get("Sentry-Hook-Resource"); resource == "" {
			slog.WarnContext(request.Context(), "[security] unauthorized request: empty Sentry-Hook-Resource")

			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if cfg.ClientToken != "" {
			body, err := io.ReadAll(request.Body)
			if err != nil {
				slog.WarnContext(request.Context(), "[security] failed to read request body")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			digest := hmac.New(sha256.New, []byte(cfg.ClientToken))

			_, err = digest.Write(body)
			if err != nil {
				slog.WarnContext(request.Context(), "[security] unauthorized request: failed to hmac")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			decodedExpectedDigest, err := hex.DecodeString(expectedDigest)
			if err != nil {
				slog.WarnContext(request.Context(), "[security] unauthorized request: failed to decode expected digest")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if !hmac.Equal(digest.Sum(nil), decodedExpectedDigest) {
				slog.WarnContext(request.Context(), "[security] unauthorized request: digest mismatch")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			slog.DebugContext(request.Context(), "[security] authorized request: request has valid digest")

			request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		slog.InfoContext(request.Context(), "[security] request authorized")

		next.ServeHTTP(w, request)
	})
}
