// Package fence provides a write-once token fence that ensures a secret
// injection window is entered and exited exactly once, preventing duplicate
// or concurrent exec runs from consuming stale environment snapshots.
package fence

import (
	"errors"
	"sync"
	"time"
)

// ErrAlreadyOpen is returned when Open is called on an already-open fence.
var ErrAlreadyOpen = errors.New("fence: already open")

// ErrNotOpen is returned when Close is called on a fence that was never opened.
var ErrNotOpen = errors.New("fence: not open")

// ErrExpired is returned when the fence TTL has elapsed before Close is called.
var ErrExpired = errors.New("fence: token expired")

// Fence guards a single injection window with an optional TTL.
type Fence struct {
	mu       sync.Mutex
	open     bool
	token    string
	openedAt time.Time
	ttl      time.Duration
	now      func() time.Time
}

// Option configures a Fence.
type Option func(*Fence)

// WithTTL sets the maximum duration the fence may remain open.
func WithTTL(d time.Duration) Option {
	return func(f *Fence) { f.ttl = d }
}

// New creates a Fence with the supplied options.
func New(opts ...Option) *Fence {
	f := &Fence{now: time.Now}
	for _, o := range opts {
		o(f)
	}
	return f
}

// Open marks the fence as active and records the provided token.
// Returns ErrAlreadyOpen if the fence is already open.
func (f *Fence) Open(token string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.open {
		return ErrAlreadyOpen
	}
	f.open = true
	f.token = token
	f.openedAt = f.now()
	return nil
}

// Close marks the fence as inactive, validating the token and TTL.
// Returns ErrNotOpen, ErrExpired, or a token mismatch error as appropriate.
func (f *Fence) Close(token string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if !f.open {
		return ErrNotOpen
	}
	if f.ttl > 0 && f.now().Sub(f.openedAt) > f.ttl {
		f.open = false
		return ErrExpired
	}
	if token != f.token {
		return errors.New("fence: token mismatch")
	}
	f.open = false
	f.token = ""
	return nil
}

// IsOpen reports whether the fence is currently open.
func (f *Fence) IsOpen() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.open
}
