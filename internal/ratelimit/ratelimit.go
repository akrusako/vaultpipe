// Package ratelimit provides a token-bucket rate limiter for controlling
// the frequency of Vault secret fetch operations.
package ratelimit

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Limiter controls the rate at which operations are allowed to proceed.
type Limiter struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	rate     float64 // tokens per second
	lastTick time.Time
	clock    func() time.Time
}

// Config holds configuration for the rate limiter.
type Config struct {
	// Rate is the number of operations allowed per second.
	Rate float64
	// Burst is the maximum number of operations allowed in a burst.
	Burst int
}

// DefaultConfig returns a Config suitable for most Vault fetch workloads.
func DefaultConfig() Config {
	return Config{
		Rate:  5,
		Burst: 10,
	}
}

// Validate checks that the Config values are usable.
func (c Config) Validate() error {
	if c.Rate <= 0 {
		return fmt.Errorf("ratelimit: rate must be greater than zero")
	}
	if c.Burst <= 0 {
		return fmt.Errorf("ratelimit: burst must be greater than zero")
	}
	return nil
}

// New creates a Limiter from the provided Config.
func New(cfg Config) (*Limiter, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Limiter{
		tokens:   float64(cfg.Burst),
		max:      float64(cfg.Burst),
		rate:     cfg.Rate,
		lastTick: time.Now(),
		clock:    time.Now,
	}, nil
}

// Wait blocks until a token is available or ctx is cancelled.
func (l *Limiter) Wait(ctx context.Context) error {
	for {
		if l.tryConsume() {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(10 * time.Millisecond):
		}
	}
}

// tryConsume attempts to consume a token. Returns true if successful.
func (l *Limiter) tryConsume() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.clock()
	elapsed := now.Sub(l.lastTick).Seconds()
	l.lastTick = now

	l.tokens += elapsed * l.rate
	if l.tokens > l.max {
		l.tokens = l.max
	}

	if l.tokens >= 1 {
		l.tokens--
		return true
	}
	return false
}
