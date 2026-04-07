package main

import (
	"fmt"
	"log"
	"net/http"

	loadbalancer "github.com/bizak0/api-gateway/internal/loadbalancer"
)

func main() {
	fmt.Println("Load Balancer starting on port 8082...")

	lb := loadbalancer.NewLoadBalancer()

	lb.Register("service-1", "http://localhost:9091")
	lb.Register("service-2", "http://localhost:9092")
	lb.Register("service-3", "http://localhost:9093")

	lb.AddRoute("/users", "users-service")
	lb.AddRoute("/orders", "orders-service")

	err := http.ListenAndServe(":8082", lb)
	if err != nil {
		log.Fatal(err)
	}
}
