// Package jitter adds randomised delay to secret-fetch intervals so that
// multiple vaultpipe processes do not stampede Vault simultaneously.
package jitter

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

// Config holds tunable parameters for the jitter spreader.
type Config struct {
	// Base is the nominal interval between fetches.
	Base time.Duration
	// Factor is the maximum fraction of Base that may be added as jitter.
	// Must be in the range (0, 1].
	Factor float64
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Base:   30 * time.Second,
		Factor: 0.25,
	}
}

func (c Config) validate() error {
	if c.Base <= 0 {
		return errors.New("jitter: Base must be positive")
	}
	if c.Factor <= 0 || c.Factor > 1 {
		return errors.New("jitter: Factor must be in range (0, 1]")
	}
	return nil
}

// Spreader sleeps for Base ± jitter before signalling readiness.
type Spreader struct {
	cfg  Config
	rand *rand.Rand
}

// New creates a Spreader from cfg. Returns an error if cfg is invalid.
func New(cfg Config) (*Spreader, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	//nolint:gosec // non-cryptographic jitter is intentional
	return &Spreader{cfg: cfg, rand: rand.New(rand.NewSource(time.Now().UnixNano()))}, nil
}

// Wait blocks for a jittered duration derived from cfg.Base. It returns early
// if ctx is cancelled, in which case ctx.Err() is returned.
func (s *Spreader) Wait(ctx context.Context) error {
	delay := s.Next()
	select {
	case <-time.After(delay):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Next returns the next jittered duration without blocking.
func (s *Spreader) Next() time.Duration {
	max := float64(s.cfg.Base) * s.cfg.Factor
	offset := time.Duration(s.rand.Float64() * max)
	return s.cfg.Base + offset
}
