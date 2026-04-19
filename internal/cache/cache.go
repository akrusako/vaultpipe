// Package cache provides a simple in-memory TTL cache for Vault secrets
// to reduce the number of requests made to the Vault server.
package cache

import (
	"sync"
	"time"
)

// Entry holds a cached secret value and its expiry time.
type Entry struct {
	Secrets   map[string]string
	FetchedAt time.Time
	TTL       time.Duration
}

// Expired returns true if the entry has passed its TTL.
func (e *Entry) Expired() bool {
	return time.Since(e.FetchedAt) > e.TTL
}

// Cache is a thread-safe in-memory store for secret entries.
type Cache struct {
	mu      sync.RWMutex
	entries map[string]*Entry
	ttl     time.Duration
}

// New creates a new Cache with the given default TTL.
func New(ttl time.Duration) *Cache {
	return &Cache{
		entries: make(map[string]*Entry),
		ttl:     ttl,
	}
}

// Set stores secrets under the given path key.
func (c *Cache) Set(path string, secrets map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[path] = &Entry{
		Secrets:   secrets,
		FetchedAt: time.Now(),
		TTL:       c.ttl,
	}
}

// Get retrieves secrets for the given path. Returns nil, false if missing or expired.
func (c *Cache) Get(path string) (map[string]string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[path]
	if !ok || e.Expired() {
		return nil, false
	}
	return e.Secrets, true
}

// Invalidate removes the entry for the given path.
func (c *Cache) Invalidate(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, path)
}

// Flush removes all entries from the cache.
func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]*Entry)
}
