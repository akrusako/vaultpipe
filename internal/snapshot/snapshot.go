// Package snapshot captures and compares secret state across rotations.
package snapshot

import (
	"sync"
	"time"
)

// Snapshot holds a point-in-time copy of secrets.
type Snapshot struct {
	Secrets   map[string]string
	CapturedAt time.Time
}

// Store maintains the latest snapshot and allows diffing against new state.
type Store struct {
	mu      sync.RWMutex
	current *Snapshot
}

// New returns an empty Store.
func New() *Store {
	return &Store{}
}

// Save records a new snapshot from the provided secrets map.
func (s *Store) Save(secrets map[string]string) {
	copy := make(map[string]string, len(secrets))
	for k, v := range secrets {
		copy[k] = v
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.current = &Snapshot{
		Secrets:    copy,
		CapturedAt: time.Now(),
	}
}

// Latest returns the most recent snapshot, or nil if none has been saved.
func (s *Store) Latest() *Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.current
}

// Changed returns true if any key differs between the stored snapshot and incoming.
func (s *Store) Changed(incoming map[string]string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.current == nil {
		return true
	}
	if len(s.current.Secrets) != len(incoming) {
		return true
	}
	for k, v := range incoming {
		if s.current.Secrets[k] != v {
			return true
		}
	}
	return false
}
