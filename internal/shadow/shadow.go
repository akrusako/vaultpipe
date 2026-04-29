// Package shadow provides a write-once secret store that prevents accidental
// overwriting of already-resolved secrets during pipeline construction.
package shadow

import (
	"errors"
	"fmt"
	"sync"
)

// ErrKeyExists is returned when an attempt is made to write a key that is
// already present in the shadow store.
var ErrKeyExists = errors.New("shadow: key already exists")

// Store is a concurrency-safe write-once map for secret key/value pairs.
type Store struct {
	mu   sync.RWMutex
	data map[string]string
}

// New returns an empty Store.
func New() *Store {
	return &Store{data: make(map[string]string)}
}

// Set writes key to the store. It returns ErrKeyExists if the key is already
// present, ensuring secrets are never silently overwritten.
func (s *Store) Set(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[key]; ok {
		return fmt.Errorf("%w: %q", ErrKeyExists, key)
	}
	s.data[key] = value
	return nil
}

// Get retrieves the value for key. The second return value reports whether
// the key was found.
func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[key]
	return v, ok
}

// Snapshot returns a shallow copy of the current store contents.
func (s *Store) Snapshot() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]string, len(s.data))
	for k, v := range s.data {
		out[k] = v
	}
	return out
}

// Len returns the number of entries currently held in the store.
func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}

// LoadMap bulk-loads secrets from m into the store, respecting the write-once
// constraint. It returns the first collision error encountered, if any.
func (s *Store) LoadMap(m map[string]string) error {
	for k, v := range m {
		if err := s.Set(k, v); err != nil {
			return err
		}
	}
	return nil
}
