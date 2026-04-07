package adaptor

import (
	"encoding/json"
	"net/http"
)

type Request struct {
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
	Body    interface{}       `json:"body"`
}

type Response struct {
	Status  int         `json:"status"`
	Headers http.Header `json:"headers"`
	Body    interface{} `json:"body"`
}

func HTTPToInternal(r *http.Request) *Request {
	headers := make(map[string]string)
	for k, v := range r.Header {
		headers[k] = v[0]
	}

	return &Request{
		Method:  r.Method,
		Path:    r.URL.Path,
		Headers: headers,
	}
}

func InternalToHTTP(w http.ResponseWriter, resp *Response) {
	for k, v := range resp.Headers {
		w.Header()[k] = v
	}

	w.WriteHeader(resp.Status)

	if resp.Body != nil {
		json.NewEncoder(w).Encode(resp.Body)
	}
}

func AdaptorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
