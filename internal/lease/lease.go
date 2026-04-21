// Package lease tracks Vault secret lease durations and emits
// renewal warnings when a lease approaches expiry.
package lease

import (
	"sync"
	"time"
)

// Info holds metadata about a single Vault secret lease.
type Info struct {
	LeaseID   string
	Path      string
	Duration  time.Duration
	Renewable bool
	Acquired  time.Time
}

// ExpiresAt returns the absolute expiry time of the lease.
func (i Info) ExpiresAt() time.Time {
	return i.Acquired.Add(i.Duration)
}

// TTL returns the remaining time-to-live relative to now.
func (i Info) TTL() time.Duration {
	return time.Until(i.ExpiresAt())
}

// Tracker stores active leases and provides lookup by lease ID.
type Tracker struct {
	mu     sync.RWMutex
	leases map[string]Info
}

// New returns an initialised Tracker.
func New() *Tracker {
	return &Tracker{
		leases: make(map[string]Info),
	}
}

// Add registers a lease. If a lease with the same ID already exists it is
// replaced.
func (t *Tracker) Add(info Info) {
	if info.Acquired.IsZero() {
		info.Acquired = time.Now()
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.leases[info.LeaseID] = info
}

// Remove deletes a lease by ID. It is a no-op if the ID is unknown.
func (t *Tracker) Remove(leaseID string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.leases, leaseID)
}

// Get returns the Info for a lease ID, and whether it was found.
func (t *Tracker) Get(leaseID string) (Info, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	info, ok := t.leases[leaseID]
	return info, ok
}

// Expiring returns all leases whose TTL is below the given threshold.
func (t *Tracker) Expiring(threshold time.Duration) []Info {
	t.mu.RLock()
	defer t.mu.RUnlock()
	var out []Info
	for _, info := range t.leases {
		if info.TTL() <= threshold {
			out = append(out, info)
		}
	}
	return out
}

// All returns a snapshot of every tracked lease.
func (t *Tracker) All() []Info {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]Info, 0, len(t.leases))
	for _, info := range t.leases {
		out = append(out, info)
	}
	return out
}
