package middleware

import (
	"net/http"
)

func RoleMiddleware(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(UserKey).(string)
			if !ok || claims == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			role := r.Header.Get("X-Role")
			if role != requiredRole {
				http.Error(w, "Forbidden - you don't have the required role: "+requiredRole, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
