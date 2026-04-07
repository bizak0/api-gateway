package proxy

import (
	"log"
	"net/http"
)

func (rp *ReverseProxy) ListenAndServeTLS(addr string, certFile string, keyFile string) {
	log.Printf("Reverse Proxy starting with TLS on %s...", addr)

	err := http.ListenAndServeTLS(addr, certFile, keyFile, rp)
	if err != nil {
		log.Fatal(err)
	}
}
