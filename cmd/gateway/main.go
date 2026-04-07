package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bizak0/api-gateway/internal/gateway/middleware"
)

func main() {
	fmt.Println("API Gateway starting on port 8081...")

	rateLimiter := middleware.NewRateLimiter(5, 10)

	mux := http.NewServeMux()

	mux.HandleFunc("/public", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Public route - no auth needed!"))
	})

	mux.Handle("/private", middleware.AuthMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := r.Context().Value(middleware.UserKey)
			w.Write([]byte("Private route - welcome " + user.(string) + "!"))
		}),
	))

	handler := rateLimiter.Middleware(mux)

	err := http.ListenAndServe(":8081", handler)
	if err != nil {
		log.Fatal(err)
	}
}
