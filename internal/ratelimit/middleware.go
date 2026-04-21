package ratelimit

import (
	"context"
	"fmt"
)

// FetchFunc is the signature of a function that fetches secrets from Vault.
type FetchFunc func(ctx context.Context, paths []string) (map[string]string, error)

// Middleware wraps a FetchFunc with rate-limiting behaviour.
type Middleware struct {
	limiter *Limiter
	next    FetchFunc
}

// NewMiddleware creates a Middleware that applies lim before every call to next.
func NewMiddleware(lim *Limiter, next FetchFunc) (*Middleware, error) {
	if lim == nil {
		return nil, fmt.Errorf("ratelimit: limiter must not be nil")
	}
	if next == nil {
		return nil, fmt.Errorf("ratelimit: next fetch func must not be nil")
	}
	return &Middleware{limiter: lim, next: next}, nil
}

// Fetch waits for a rate-limit token then delegates to the wrapped FetchFunc.
func (m *Middleware) Fetch(ctx context.Context, paths []string) (map[string]string, error) {
	if err := m.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("ratelimit: wait: %w", err)
	}
	return m.next(ctx, paths)
}
