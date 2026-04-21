package throttle

import (
	"context"
	"fmt"
)

// Func is a generic operation that returns a map of secrets.
type Func func(ctx context.Context, paths []string) (map[string]string, error)

// Middleware wraps a Func with concurrency limiting so that at most
// Limiter.Concurrency fetch calls run simultaneously.
type Middleware struct {
	limiter *Limiter
	next    Func
}

// NewMiddleware creates a Middleware that gates next with l.
func NewMiddleware(l *Limiter, next Func) (*Middleware, error) {
	if l == nil {
		return nil, fmt.Errorf("throttle: limiter must not be nil")
	}
	if next == nil {
		return nil, fmt.Errorf("throttle: next func must not be nil")
	}
	return &Middleware{limiter: l, next: next}, nil
}

// Fetch acquires a slot, delegates to the wrapped Func, then releases the slot.
func (m *Middleware) Fetch(ctx context.Context, paths []string) (map[string]string, error) {
	if err := m.limiter.Acquire(ctx); err != nil {
		return nil, fmt.Errorf("throttle middleware: %w", err)
	}
	defer m.limiter.Release()
	return m.next(ctx, paths)
}
