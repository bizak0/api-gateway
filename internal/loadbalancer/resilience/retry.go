package resilience

import (
	"log"
	"net/http"
	"time"
)

type RetryConfig struct {
	MaxRetries int
	Delay      time.Duration
}

func NewRetryConfig(maxRetries int, delay time.Duration) *RetryConfig {
	return &RetryConfig{
		MaxRetries: maxRetries,
		Delay:      delay,
	}
}

func (rc *RetryConfig) Do(req *http.Request) (*http.Response, error) {
	client := &http.Client{Timeout: 5 * time.Second}

	var lastErr error
	for i := 0; i <= rc.MaxRetries; i++ {
		if i > 0 {
			log.Printf("Retry %d/%d for %s", i, rc.MaxRetries, req.URL.Path)
			time.Sleep(rc.Delay)
		}

		resp, err := client.Do(req)
		if err == nil {
			return resp, nil
		}

		lastErr = err
	}

	return nil, lastErr
}
