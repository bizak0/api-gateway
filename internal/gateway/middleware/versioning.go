package middleware

import (
	"net/http"
	"strings"
)

func VersionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if strings.HasPrefix(path, "/v1") {
			w.Header().Set("X-API-Version", "v1")
			next.ServeHTTP(w, r)
			return
		}

		if strings.HasPrefix(path, "/v2") {
			w.Header().Set("X-API-Version", "v2")
			next.ServeHTTP(w, r)
			return
		}

		http.Error(w, "API version not specified. Use /v1 or /v2", http.StatusBadRequest)
	})
}
