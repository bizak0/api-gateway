package cache

import (
	"net/http"
	"sync"
	"time"
)

type CacheEntry struct {
	Body      []byte
	Headers   http.Header
	Timestamp time.Time
}

type PublicCache struct {
	mu      sync.Mutex
	entries map[string]CacheEntry
	ttl     time.Duration
}

func NewPublicCache(ttl time.Duration) *PublicCache {
	return &PublicCache{
		entries: make(map[string]CacheEntry),
		ttl:     ttl,
	}
}

func (c *PublicCache) Get(key string) (CacheEntry, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.entries[key]
	if !exists {
		return CacheEntry{}, false
	}

	if time.Since(entry.Timestamp) > c.ttl {
		delete(c.entries, key)
		return CacheEntry{}, false
	}

	return entry, true
}

func (c *PublicCache) Set(key string, body []byte, headers http.Header) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = CacheEntry{
		Body:      body,
		Headers:   headers,
		Timestamp: time.Now(),
	}
}
