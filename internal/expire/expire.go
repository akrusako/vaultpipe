// Package expire provides TTL-based expiration tracking for secrets and
// arbitrary key-value pairs. It allows callers to register items with a
// time-to-live and query whether they have expired.
package expire

import (
	"sync"
	"time"
)

// Entry holds an item and its expiration deadline.
type Entry struct {
	Value     string
	ExpiresAt time.Time
}

// IsExpired reports whether the entry has passed its deadline.
func (e Entry) IsExpired(now time.Time) bool {
	return now.After(e.ExpiresAt)
}

// Tracker manages TTL expiration for named entries.
type Tracker struct {
	mu      sync.RWMutex
	entries map[string]Entry
	now     func() time.Time
}

// New returns a new Tracker. If now is nil, time.Now is used.
func New(now func() time.Time) *Tracker {
	if now == nil {
		now = time.Now
	}
	return &Tracker{
		entries: make(map[string]Entry),
		now:     now,
	}
}

// Set registers key with value and a TTL duration.
func (t *Tracker) Set(key, value string, ttl time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries[key] = Entry{
		Value:     value,
		ExpiresAt: t.now().Add(ttl),
	}
}

// Get returns the entry for key and whether it exists and has not expired.
func (t *Tracker) Get(key string) (Entry, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	e, ok := t.entries[key]
	if !ok {
		return Entry{}, false
	}
	if e.IsExpired(t.now()) {
		return Entry{}, false
	}
	return e, true
}

// Delete removes the entry for key.
func (t *Tracker) Delete(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, key)
}

// Purge removes all expired entries and returns the count removed.
func (t *Tracker) Purge() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.now()
	removed := 0
	for k, e := range t.entries {
		if e.IsExpired(now) {
			delete(t.entries, k)
			removed++
		}
	}
	return removed
}

// Len returns the total number of tracked entries, including expired ones.
func (t *Tracker) Len() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.entries)
}
