// Package backoff provides configurable exponential backoff with jitter
// for use when retrying transient failures against Vault or child processes.
package backoff

import (
	"errors"
	"math"
	"math/rand"
	"time"
)

// Config holds parameters that control the backoff behaviour.
type Config struct {
	// InitialInterval is the delay before the first retry.
	InitialInterval time.Duration
	// MaxInterval caps the computed delay regardless of attempt count.
	MaxInterval time.Duration
	// Multiplier is applied to the previous interval on each attempt.
	Multiplier float64
	// Jitter adds a random fraction of the computed interval to avoid
	// thundering-herd scenarios. Values in [0, 1].
	Jitter float64
}

// DefaultConfig returns sensible defaults suitable for Vault API calls.
func DefaultConfig() Config {
	return Config{
		InitialInterval: 500 * time.Millisecond,
		MaxInterval:     30 * time.Second,
		Multiplier:      1.5,
		Jitter:          0.2,
	}
}

// Validate returns an error if any field is out of range.
func (c Config) Validate() error {
	if c.InitialInterval <= 0 {
		return errors.New("backoff: InitialInterval must be positive")
	}
	if c.MaxInterval < c.InitialInterval {
		return errors.New("backoff: MaxInterval must be >= InitialInterval")
	}
	if c.Multiplier < 1.0 {
		return errors.New("backoff: Multiplier must be >= 1.0")
	}
	if c.Jitter < 0 || c.Jitter > 1 {
		return errors.New("backoff: Jitter must be in [0, 1]")
	}
	return nil
}

// Backoff computes per-attempt delays.
type Backoff struct {
	cfg     Config
	attempt int
}

// New creates a Backoff from cfg. Returns an error if cfg is invalid.
func New(cfg Config) (*Backoff, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Backoff{cfg: cfg}, nil
}

// Next returns the delay to wait before the next attempt and advances the
// internal counter. The returned duration is always in
// [InitialInterval, MaxInterval].
func (b *Backoff) Next() time.Duration {
	raw := float64(b.cfg.InitialInterval) * math.Pow(b.cfg.Multiplier, float64(b.attempt))
	if raw > float64(b.cfg.MaxInterval) {
		raw = float64(b.cfg.MaxInterval)
	}
	if b.cfg.Jitter > 0 {
		raw += raw * b.cfg.Jitter * rand.Float64() //nolint:gosec
		if raw > float64(b.cfg.MaxInterval) {
			raw = float64(b.cfg.MaxInterval)
		}
	}
	b.attempt++
	return time.Duration(raw)
}

// Reset sets the attempt counter back to zero.
func (b *Backoff) Reset() {
	b.attempt = 0
}

// Attempt returns the current (zero-based) attempt index.
func (b *Backoff) Attempt() int {
	return b.attempt
}
