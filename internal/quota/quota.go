// Package quota enforces per-path secret read limits within a sliding window.
package quota

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// ErrQuotaExceeded is returned when a path has exceeded its allowed reads.
var ErrQuotaExceeded = errors.New("quota: read limit exceeded for path")

// Config holds quota configuration.
type Config struct {
	// MaxReads is the maximum number of reads allowed per window.
	MaxReads int
	// Window is the duration of the sliding window.
	Window time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		MaxReads: 100,
		Window:   time.Minute,
	}
}

func (c Config) validate() error {
	if c.MaxReads <= 0 {
		return errors.New("quota: MaxReads must be greater than zero")
	}
	if c.Window <= 0 {
		return errors.New("quota: Window must be greater than zero")
	}
	return nil
}

type entry struct {
	timestamps []time.Time
}

// Quota tracks and enforces read limits per secret path.
type Quota struct {
	cfg Config
	mu  sync.Mutex
	map_ map[string]*entry
	now  func() time.Time
}

// New creates a new Quota enforcer with the given config.
func New(cfg Config) (*Quota, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return &Quota{
		cfg:  cfg,
		map_: make(map[string]*entry),
		now:  time.Now,
	}, nil
}

// Check returns an error if path has exceeded its quota, otherwise records the read.
func (q *Quota) Check(path string) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	now := q.now()
	cutoff := now.Add(-q.cfg.Window)

	e, ok := q.map_[path]
	if !ok {
		e = &entry{}
		q.map_[path] = e
	}

	// Evict timestamps outside the window.
	valid := e.timestamps[:0]
	for _, t := range e.timestamps {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	e.timestamps = valid

	if len(e.timestamps) >= q.cfg.MaxReads {
		return fmt.Errorf("%w: %s", ErrQuotaExceeded, path)
	}

	e.timestamps = append(e.timestamps, now)
	return nil
}

// Reset clears all recorded reads for the given path.
func (q *Quota) Reset(path string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	delete(q.map_, path)
}

// Counts returns the current read count within the window for each tracked path.
func (q *Quota) Counts() map[string]int {
	q.mu.Lock()
	defer q.mu.Unlock()

	now := q.now()
	cutoff := now.Add(-q.cfg.Window)
	out := make(map[string]int, len(q.map_))
	for path, e := range q.map_ {
		count := 0
		for _, t := range e.timestamps {
			if t.After(cutoff) {
				count++
			}
		}
		if count > 0 {
			out[path] = count
		}
	}
	return out
}
