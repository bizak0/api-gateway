package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/bizak0/api-gateway/internal/gateway/adaptor"
	"github.com/bizak0/api-gateway/internal/gateway/middleware"
)

func main() {
	fmt.Println("API Gateway starting on port 8081...")
	fmt.Println("Forwarding to Load Balancer on port 8082...")

	rateLimiter := middleware.NewRateLimiter(5, 10)

	// ✅ NOUVEAU : reverse proxy interne vers le Load Balancer
	lbURL, err := url.Parse("http://localhost:8082")
	if err != nil {
		log.Fatal("Erreur parsing URL Load Balancer:", err)
	}
	lbProxy := httputil.NewSingleHostReverseProxy(lbURL)

	mux := http.NewServeMux()

	// Route publique → middleware adaptor → forward vers LB
	mux.HandleFunc("/public", func(w http.ResponseWriter, r *http.Request) {
		adaptor.HTTPToInternal(r) // transformation interne conservée
		lbProxy.ServeHTTP(w, r)
	})

	// Route privée → Auth requis → forward vers LB
	mux.Handle("/private", middleware.AuthMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lbProxy.ServeHTTP(w, r)
		}),
	))

	// Route admin → Auth + Role admin requis → forward vers LB
	mux.Handle("/admin", middleware.AuthMiddleware(
		middleware.RoleMiddleware("admin")(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				lbProxy.ServeHTTP(w, r)
			}),
		),
	))

	// Routes versionnées → forward vers LB
	mux.HandleFunc("/v1/users", func(w http.ResponseWriter, r *http.Request) {
		lbProxy.ServeHTTP(w, r)
	})

	mux.HandleFunc("/v2/users", func(w http.ResponseWriter, r *http.Request) {
		lbProxy.ServeHTTP(w, r)
	})

	// ✅ Route catch-all : toute autre requête → forward vers LB
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lbProxy.ServeHTTP(w, r)
	})

	// Chaîne de middleware conservée intégralement
	handler := rateLimiter.Middleware(
		middleware.TransformMiddleware(
			adaptor.AdaptorMiddleware(
				middleware.VersionMiddleware(mux),
			),
		),
	)

	if err := http.ListenAndServe(":8081", handler); err != nil {
		log.Fatal(err)
	}
}
