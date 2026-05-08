package main

import (
	"fmt"

	"github.com/bizak0/api-gateway/internal/proxy"
)

func main() {
	fmt.Println("Reverse Proxy starting with TLS on port 8443...")
	fmt.Println("Forwarding to API Gateway on port 8081...")

	rp := proxy.NewReverseProxy("http://localhost:8081")
	rp.ListenAndServeTLS(":8443", "cert.pem", "key.pem")
}
