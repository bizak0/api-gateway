package main

import (
	"fmt"
	"log"
	"net/http"

	loadbalancer "github.com/bizak0/api-gateway/internal/loadbalancer"
)

func main() {
	fmt.Println("Load Balancer starting on port 8082...")
	fmt.Println("Distributing requests to microservices...")

	lb := loadbalancer.NewLoadBalancer()

	// ✅ CORRECTION : noms de services cohérents avec les routes déclarées
	lb.Register("users-service", "http://localhost:9091")
	lb.Register("orders-service", "http://localhost:9092")
	lb.Register("admin-service", "http://localhost:9093")

	// ✅ CORRECTION : routes alignées avec celles du Gateway
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
