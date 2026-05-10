package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	loadbalancer "github.com/bizak0/api-gateway/internal/loadbalancer"
)

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func main() {
	fmt.Println("Load Balancer starting on port 8082...")

	usersURL := getEnv("USERS_SERVICE_URL", "http://localhost:9091")
	ordersURL := getEnv("ORDERS_SERVICE_URL", "http://localhost:9092")
	adminURL := getEnv("ADMIN_SERVICE_URL", "http://localhost:9093")

	lb := loadbalancer.NewLoadBalancer(
		10*time.Second,
		3,
		30*time.Second,
		3,
		1*time.Second,
	)

	lb.Register("users-service", usersURL)
	lb.Register("orders-service", ordersURL)
	lb.Register("admin-service", adminURL)

	lb.AddRoute("/public", "users-service")
	lb.AddRoute("/private", "users-service")
	lb.AddRoute("/v1/users", "users-service")
	lb.AddRoute("/v2/users", "users-service")
	lb.AddRoute("/orders", "orders-service")
	lb.AddRoute("/admin", "admin-service")

	if err := http.ListenAndServe(":8082", lb); err != nil {
		log.Fatal(err)
	}
}
