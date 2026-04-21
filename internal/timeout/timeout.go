// Package timeout provides context-based deadline enforcement for
// child process execution and secret fetch operations.
package timeout

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ErrDeadlineExceeded is returned when an operation exceeds its configured timeout.
var ErrDeadlineExceeded = errors.New("timeout: deadline exceeded")

// Config holds timeout configuration.
type Config struct {
	// Exec is the maximum duration allowed for the child process.
	Exec time.Duration
	// Fetch is the maximum duration allowed for a single secret fetch.
	Fetch time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Exec:  0, // 0 means no timeout
		Fetch: 10 * time.Second,
	}
}

// Validate returns an error if the Config is invalid.
func (c Config) Validate() error {
	if c.Fetch < 0 {
		return fmt.Errorf("timeout: fetch duration must be non-negative, got %s", c.Fetch)
	}
	if c.Exec < 0 {
		return fmt.Errorf("timeout: exec duration must be non-negative, got %s", c.Exec)
	}
	return nil
}

// Enforcer applies timeouts to contexts.
type Enforcer struct {
	cfg Config
}

// New creates a new Enforcer from the given Config.
func New(cfg Config) (*Enforcer, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Enforcer{cfg: cfg}, nil
}

// ExecContext returns a child context that is cancelled after the configured
// exec timeout. If the timeout is zero the parent context is returned unchanged
// along with a no-op cancel function.
func (e *Enforcer) ExecContext(parent context.Context) (context.Context, context.CancelFunc) {
	if e.cfg.Exec <= 0 {
		return parent, func() {}
	}
	return context.WithTimeout(parent, e.cfg.Exec)
}

// FetchContext returns a child context that is cancelled after the configured
// fetch timeout.
func (e *Enforcer) FetchContext(parent context.Context) (context.Context, context.CancelFunc) {
	if e.cfg.Fetch <= 0 {
		return parent, func() {}
	}
	return context.WithTimeout(parent, e.cfg.Fetch)
}

// IsDeadline reports whether err represents a deadline/timeout condition.
func IsDeadline(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, ErrDeadlineExceeded) {
		return true
	}
	return false
}
