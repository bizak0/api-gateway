package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func TransformMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		r.Header.Set("X-Request-ID", uuid.New().String())
		r.Header.Set("X-Forwarded-By", "API-Gateway")

		w.Header().Set("X-Powered-By", "API-Gateway")
		w.Header().Set("X-API-Gateway", "true")

		next.ServeHTTP(w, r)

		duration := time.Since(start).Milliseconds()
		w.Header().Set("X-Response-Time", strconv.FormatInt(duration, 10)+"ms")
	})
}
