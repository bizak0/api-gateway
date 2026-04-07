package main

import (
	"fmt"

	"github.com/bizak0/api-gateway/internal/proxy"
)

func main() {
	fmt.Println("Reverse Proxy starting with TLS on port 8443...")

	rp := proxy.NewReverseProxy("http://localhost:9090")

	rp.ListenAndServeTLS(":8443", "cert.pem", "key.pem")
}
