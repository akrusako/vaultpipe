// Package truncate provides utilities for truncating secret values
// to a maximum byte length before they are injected into the environment.
// This prevents oversized values from exceeding OS environment variable limits.
package truncate

import "fmt"

// DefaultMaxBytes is the default maximum byte length for a secret value.
const DefaultMaxBytes = 32 * 1024 // 32 KiB

// Truncator trims secret values that exceed a configured byte limit.
type Truncator struct {
	maxBytes int
}

// Option configures a Truncator.
type Option func(*Truncator)

// WithMaxBytes sets the maximum allowed byte length for a value.
// Panics if max is less than or equal to zero.
func WithMaxBytes(max int) Option {
	if max <= 0 {
		panic(fmt.Sprintf("truncate: maxBytes must be positive, got %d", max))
	}
	return func(t *Truncator) {
		t.maxBytes = max
	}
}

// New creates a Truncator with optional configuration.
func New(opts ...Option) *Truncator {
	t := &Truncator{maxBytes: DefaultMaxBytes}
	for _, o := range opts {
		o(t)
	}
	return t
}

// Apply iterates over the provided secret map and truncates any value
// whose byte length exceeds the configured maximum. The original map
// is not modified; a new map is returned.
func (t *Truncator) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if len(v) > t.maxBytes {
			v = v[:t.maxBytes]
		}
		out[k] = v
	}
	return out
}

// MaxBytes returns the configured maximum byte length.
func (t *Truncator) MaxBytes() int {
	return t.maxBytes
}
