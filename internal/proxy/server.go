package proxy

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/bizak0/api-gateway/internal/proxy/cache"
)

type ReverseProxy struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
	cache  *cache.PublicCache
}

func NewReverseProxy(targetURL string) *ReverseProxy {
	target, err := url.Parse(targetURL)
	if err != nil {
		log.Fatal(err)
	}

	return &ReverseProxy{
		target: target,
		proxy:  httputil.NewSingleHostReverseProxy(target),
		cache:  cache.NewPublicCache(5 * time.Minute),
	}
}

func (rp *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cacheKey := r.URL.Path

	if entry, found := rp.cache.Get(cacheKey); found {
		log.Printf("Cache HIT: %s", cacheKey)
		for k, v := range entry.Headers {
			w.Header()[k] = v
		}
		w.Write(entry.Body)
		return
	}

	log.Printf("Cache MISS: %s %s", r.Method, r.URL.Path)

	rec := &responseRecorder{
		ResponseWriter: w,
		body:           &bytes.Buffer{},
	}

	rp.proxy.ServeHTTP(rec, r)
	rp.cache.Set(cacheKey, rec.body.Bytes(), rec.Header())
}

type responseRecorder struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (r *responseRecorder) Body() []byte {
	return r.body.Bytes()
}

func (r *responseRecorder) Headers() http.Header {
	_ = io.Discard
	return r.ResponseWriter.Header()
}
