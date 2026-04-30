// Package ceiling enforces an upper bound on the number of secrets fetched
// in a single request, preventing runaway reads from large Vault mounts.
package ceiling

import (
	"context"
	"fmt"
)

// DefaultMax is the default maximum number of secret keys allowed.
const DefaultMax = 256

// Fetcher is the interface satisfied by any secret-fetching layer.
type Fetcher interface {
	Fetch(ctx context.Context, paths []string) (map[string]string, error)
}

// FetchFunc is a convenience adapter for plain functions.
type FetchFunc func(ctx context.Context, paths []string) (map[string]string, error)

func (f FetchFunc) Fetch(ctx context.Context, paths []string) (map[string]string, error) {
	return f(ctx, paths)
}

// Limiter wraps a Fetcher and returns an error when the result exceeds Max keys.
type Limiter struct {
	next Fetcher
	max  int
}

// New returns a Limiter that rejects responses larger than max keys.
// If max is <= 0, DefaultMax is used.
func New(next Fetcher, max int) *Limiter {
	if max <= 0 {
		max = DefaultMax
	}
	return &Limiter{next: next, max: max}
}

// Fetch delegates to the wrapped Fetcher and enforces the key ceiling.
func (l *Limiter) Fetch(ctx context.Context, paths []string) (map[string]string, error) {
	secrets, err := l.next.Fetch(ctx, paths)
	if err != nil {
		return nil, err
	}
	if len(secrets) > l.max {
		return nil, fmt.Errorf("ceiling: response contains %d keys, limit is %d", len(secrets), l.max)
	}
	return secrets, nil
}

// Max returns the configured key ceiling.
func (l *Limiter) Max() int { return l.max }
