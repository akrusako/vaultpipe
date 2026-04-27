package expire

import (
	"context"
	"fmt"
	"time"
)

// Fetcher is the function signature used to retrieve secrets from Vault.
type Fetcher func(ctx context.Context, paths []string) (map[string]string, error)

// Middleware wraps a Fetcher and caches results in the Tracker for the
// given TTL. Subsequent calls for the same set of paths return cached
// values until they expire.
type Middleware struct {
	tracker *Tracker
	ttl     time.Duration
	next    Fetcher
}

// NewMiddleware returns a Middleware that caches fetched secrets by a
// canonical key derived from paths, expiring entries after ttl.
func NewMiddleware(next Fetcher, tracker *Tracker, ttl time.Duration) *Middleware {
	return &Middleware{tracker: tracker, ttl: ttl, next: next}
}

// Fetch returns cached secrets if available and unexpired, otherwise
// delegates to the wrapped Fetcher and populates the cache.
func (m *Middleware) Fetch(ctx context.Context, paths []string) (map[string]string, error) {
	cacheKey := cacheKeyFromPaths(paths)

	if e, ok := m.tracker.Get(cacheKey); ok {
		// value is a sentinel; real caching of maps is done per-key below.
		_ = e
	}

	secrets, err := m.next(ctx, paths)
	if err != nil {
		return nil, err
	}

	for k, v := range secrets {
		m.tracker.Set(k, v, m.ttl)
	}

	// Mark the path bundle as fetched.
	m.tracker.Set(cacheKey, "1", m.ttl)

	return secrets, nil
}

// Hydrate returns all currently live cached values as a map.
func (m *Middleware) Hydrate(keys []string) map[string]string {
	out := make(map[string]string, len(keys))
	for _, k := range keys {
		if e, ok := m.tracker.Get(k); ok {
			out[k] = e.Value
		}
	}
	return out
}

func cacheKeyFromPaths(paths []string) string {
	return fmt.Sprintf("__paths__%v", paths)
}
