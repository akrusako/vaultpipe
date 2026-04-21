// Package throttle provides a concurrency limiter that caps the number of
// simultaneous secret-fetch operations issued against Vault.
package throttle

import (
	"context"
	"errors"
	"fmt"
)

// ErrThrottled is returned when the limiter cannot accept a new ticket because
// the context has been cancelled or the semaphore is full and the caller chose
// a non-blocking acquire.
var ErrThrottled = errors.New("throttle: request rejected")

// Config holds tunables for the Limiter.
type Config struct {
	// Concurrency is the maximum number of in-flight operations.
	Concurrency int
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{Concurrency: 8}
}

// Validate returns an error when the Config is unusable.
func (c Config) Validate() error {
	if c.Concurrency <= 0 {
		return fmt.Errorf("throttle: concurrency must be > 0, got %d", c.Concurrency)
	}
	return nil
}

// Limiter is a semaphore-backed concurrency gate.
type Limiter struct {
	sem chan struct{}
}

// New creates a Limiter from cfg. Returns an error when cfg is invalid.
func New(cfg Config) (*Limiter, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Limiter{sem: make(chan struct{}, cfg.Concurrency)}, nil
}

// Acquire blocks until a slot is available or ctx is done.
// Callers must call Release exactly once after the guarded work completes.
func (l *Limiter) Acquire(ctx context.Context) error {
	select {
	case l.sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("%w: %w", ErrThrottled, ctx.Err())
	}
}

// Release frees a previously acquired slot.
func (l *Limiter) Release() {
	<-l.sem
}

// Inflight returns the number of currently held slots.
func (l *Limiter) Inflight() int {
	return len(l.sem)
}
