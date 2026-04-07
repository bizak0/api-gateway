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

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API Gateway is running!"))
	})

	handler := rateLimiter.Middleware(mux)

	err := http.ListenAndServe(":8081", handler)
	if err != nil {
		log.Fatal(err)
	}
}
