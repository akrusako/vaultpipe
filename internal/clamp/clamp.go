// Package clamp provides a middleware that enforces minimum and maximum
// byte-length constraints on secret values, clamping or rejecting values
// that fall outside the configured bounds.
package clamp

import (
	"context"
	"fmt"
)

// Fetcher retrieves a map of secret key/value pairs.
type Fetcher func(ctx context.Context, paths []string) (map[string]string, error)

// Options configures the Clamp middleware.
type Options struct {
	// MinBytes is the minimum allowed byte length for a value (0 = no minimum).
	MinBytes int
	// MaxBytes is the maximum allowed byte length for a value (0 = no maximum).
	MaxBytes int
	// DropShort silently drops values shorter than MinBytes instead of erroring.
	DropShort bool
	// DropLong silently drops values longer than MaxBytes instead of erroring.
	DropLong bool
}

// Clamp enforces byte-length constraints on fetched secret values.
type Clamp struct {
	opts Options
}

// New returns a Clamp configured with the given options.
// It returns an error if MaxBytes is set and less than MinBytes.
func New(opts Options) (*Clamp, error) {
	if opts.MaxBytes > 0 && opts.MinBytes > 0 && opts.MaxBytes < opts.MinBytes {
		return nil, fmt.Errorf("clamp: MaxBytes (%d) must be >= MinBytes (%d)", opts.MaxBytes, opts.MinBytes)
	}
	return &Clamp{opts: opts}, nil
}

// Apply filters the provided secrets map, enforcing min/max byte-length
// constraints. Values that violate a bound are either dropped (if the
// corresponding Drop flag is set) or cause Apply to return an error.
func (c *Clamp) Apply(secrets map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		n := len(v)
		if c.opts.MinBytes > 0 && n < c.opts.MinBytes {
			if c.opts.DropShort {
				continue
			}
			return nil, fmt.Errorf("clamp: value for key %q is %d bytes, minimum is %d", k, n, c.opts.MinBytes)
		}
		if c.opts.MaxBytes > 0 && n > c.opts.MaxBytes {
			if c.opts.DropLong {
				continue
			}
			return nil, fmt.Errorf("clamp: value for key %q is %d bytes, maximum is %d", k, n, c.opts.MaxBytes)
		}
		out[k] = v
	}
	return out, nil
}

// NewMiddleware wraps next with Clamp enforcement. Secrets returned by next
// are passed through Apply before being returned to the caller.
func NewMiddleware(c *Clamp, next Fetcher) Fetcher {
	return func(ctx context.Context, paths []string) (map[string]string, error) {
		secrets, err := next(ctx, paths)
		if err != nil {
			return nil, err
		}
		return c.Apply(secrets)
	}
}
