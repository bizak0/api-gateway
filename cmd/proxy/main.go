package main

import (
	"fmt"
	"os"

	"github.com/bizak0/api-gateway/internal/proxy"
)

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func main() {
	gatewayURL := getEnv("GATEWAY_URL", "http://localhost:8081")

	fmt.Println("Reverse Proxy starting with TLS on port 8443...")
	fmt.Printf("Forwarding to API Gateway at %s...\n", gatewayURL)

	rp := proxy.NewReverseProxy(gatewayURL)
	rp.ListenAndServeTLS(":8443", "cert.pem", "key.pem")
}
