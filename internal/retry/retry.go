// Package retry provides simple exponential backoff retry logic
// for transient errors encountered during secret retrieval or token renewal.
package retry

import (
	"context"
	"errors"
	"time"
)

// ErrMaxAttempts is returned when all retry attempts are exhausted.
var ErrMaxAttempts = errors.New("retry: max attempts reached")

// Config holds retry parameters.
type Config struct {
	MaxAttempts int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
}

// DefaultConfig returns sensible retry defaults.
func DefaultConfig() Config {
	return Config{
		MaxAttempts:  5,
		InitialDelay: 200 * time.Millisecond,
		MaxDelay:     10 * time.Second,
		Multiplier:   2.0,
	}
}

// Do calls fn up to cfg.MaxAttempts times, backing off exponentially.
// It stops early if ctx is cancelled or fn returns a non-retryable error.
func Do(ctx context.Context, cfg Config, fn func() error) error {
	delay := cfg.InitialDelay
	var lastErr error
	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		var nr *NonRetryableError
		if errors.As(lastErr, &nr) {
			return nr.Unwrap()
		}
		if attempt < cfg.MaxAttempts-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
			delay = time.Duration(float64(delay) * cfg.Multiplier)
			if delay > cfg.MaxDelay {
				delay = cfg.MaxDelay
			}
		}
	}
	return ErrMaxAttempts
}

// NonRetryableError wraps an error to signal that retry should stop immediately.
type NonRetryableError struct {
	cause error
}

// Permanent wraps err so that Do will not retry.
func Permanent(err error) error {
	if err == nil {
		return nil
	}
	return &NonRetryableError{cause: err}
}

func (e *NonRetryableError) Error() string { return e.cause.Error() }
func (e *NonRetryableError) Unwrap() error { return e.cause }
