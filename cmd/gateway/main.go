package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bizak0/api-gateway/internal/gateway/adaptor"
	"github.com/bizak0/api-gateway/internal/gateway/middleware"
)

func main() {
	fmt.Println("API Gateway starting on port 8081...")

	rateLimiter := middleware.NewRateLimiter(5, 10)

	mux := http.NewServeMux()

	mux.HandleFunc("/public", func(w http.ResponseWriter, r *http.Request) {
		req := adaptor.HTTPToInternal(r)
		w.Write([]byte("Public route - Method: " + req.Method + " Path: " + req.Path))
	})

	mux.Handle("/private", middleware.AuthMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := r.Context().Value(middleware.UserKey)
			w.Write([]byte("Private route - welcome " + user.(string) + "!"))
		}),
	))

	mux.Handle("/admin", middleware.AuthMiddleware(
		middleware.RoleMiddleware("admin")(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Admin route - welcome admin!"))
			}),
		),
	))

	mux.HandleFunc("/v1/users", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("V1 - Users list"))
	})

	mux.HandleFunc("/v2/users", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("V2 - Users list with more details!"))
	})

	handler := rateLimiter.Middleware(
		middleware.TransformMiddleware(
			adaptor.AdaptorMiddleware(
				middleware.VersionMiddleware(mux),
			),
		),
	)

	err := http.ListenAndServe(":8081", handler)
	if err != nil {
		log.Fatal(err)
	}
}
