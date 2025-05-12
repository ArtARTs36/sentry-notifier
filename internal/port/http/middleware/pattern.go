package middleware

import (
	"net/http"
	"strings"
)

func Pattern(pattern string, target http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if !strings.HasPrefix(req.RequestURI, pattern) {
			http.NotFound(w, req)
			return
		}

		target.ServeHTTP(w, req)
	})
}
