package cache

import (
	"net/http"
	"sync"
	"time"
)

type PrivateCacheEntry struct {
	Body      []byte
	Headers   http.Header
	Timestamp time.Time
}

type PrivateCache struct {
	mu      sync.Mutex
	entries map[string]PrivateCacheEntry
	ttl     time.Duration
}

func NewPrivateCache(ttl time.Duration) *PrivateCache {
	return &PrivateCache{
		entries: make(map[string]PrivateCacheEntry),
		ttl:     ttl,
	}
}

func (c *PrivateCache) Get(userID string, path string) (PrivateCacheEntry, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := userID + ":" + path
	entry, exists := c.entries[key]
	if !exists {
		return PrivateCacheEntry{}, false
	}

	if time.Since(entry.Timestamp) > c.ttl {
		delete(c.entries, key)
		return PrivateCacheEntry{}, false
	}

	return entry, true
}

func (c *PrivateCache) Set(userID string, path string, body []byte, headers http.Header) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := userID + ":" + path
	c.entries[key] = PrivateCacheEntry{
		Body:      body,
		Headers:   headers,
		Timestamp: time.Now(),
	}
}
