// Package circuit implements a circuit breaker to protect against
// repeated failures when communicating with Vault or other external services.
package circuit

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// State represents the current state of the circuit breaker.
type State int

const (
	StateClosed   State = iota // normal operation
	StateOpen                  // failing, requests blocked
	StateHalfOpen              // testing if service recovered
)

// ErrOpen is returned when the circuit breaker is open.
var ErrOpen = errors.New("circuit breaker is open")

// Config holds configuration for the circuit breaker.
type Config struct {
	// MaxFailures is the number of consecutive failures before opening.
	MaxFailures int
	// ResetTimeout is how long to wait before moving to half-open.
	ResetTimeout time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		MaxFailures:  5,
		ResetTimeout: 30 * time.Second,
	}
}

// Breaker is a circuit breaker that tracks consecutive failures.
type Breaker struct {
	mu           sync.Mutex
	cfg          Config
	state        State
	failures     int
	lastFailure  time.Time
}

// New creates a new Breaker with the given configuration.
func New(cfg Config) (*Breaker, error) {
	if cfg.MaxFailures <= 0 {
		return nil, fmt.Errorf("circuit: MaxFailures must be > 0, got %d", cfg.MaxFailures)
	}
	if cfg.ResetTimeout <= 0 {
		return nil, fmt.Errorf("circuit: ResetTimeout must be > 0, got %v", cfg.ResetTimeout)
	}
	return &Breaker{cfg: cfg, state: StateClosed}, nil
}

// Allow reports whether the call should be allowed through.
// It transitions an open breaker to half-open if the reset timeout has elapsed.
func (b *Breaker) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(b.lastFailure) >= b.cfg.ResetTimeout {
			b.state = StateHalfOpen
			return true
		}
		return false
	case StateHalfOpen:
		return true
	}
	return false
}

// RecordSuccess records a successful call, resetting the breaker to closed.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = StateClosed
}

// RecordFailure records a failed call and may open the breaker.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	b.lastFailure = time.Now()
	if b.failures >= b.cfg.MaxFailures {
		b.state = StateOpen
	}
}

// State returns the current state of the breaker.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
