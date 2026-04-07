package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bizak0/api-gateway/internal/proxy"
)

func main() {
	fmt.Println("Reverse Proxy starting on port 8080...")

	rp := proxy.NewReverseProxy("http://localhost:9090")

	err := http.ListenAndServe(":8080", rp)
	if err != nil {
		log.Fatal(err)
	}
}
