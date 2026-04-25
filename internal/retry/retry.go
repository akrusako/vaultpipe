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
// If all attempts are exhausted, ErrMaxAttempts is returned wrapping the
// last error encountered, so callers can inspect the underlying cause via
// errors.Unwrap or errors.Is.
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
	return &maxAttemptsError{cause: lastErr}
}

// maxAttemptsError is returned when all retry attempts are exhausted,
// wrapping the last error so callers can inspect the root cause.
type maxAttemptsError struct {
	cause error
}

func (e *maxAttemptsError) Error() string {
	if e.cause != nil {
		return ErrMaxAttempts.Error() + ": " + e.cause.Error()
	}
	return ErrMaxAttempts.Error()
}
func (e *maxAttemptsError) Unwrap() error { return e.cause }
func (e *maxAttemptsError) Is(target error) bool { return target == ErrMaxAttempts }

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
