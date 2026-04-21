// Package metrics provides lightweight runtime counters for vaultpipe
// operations such as secret fetches, cache hits, and exec invocations.
package metrics

import (
	"sync"
	"sync/atomic"
)

// Counter is a monotonically increasing uint64 counter.
type Counter struct{ v uint64 }

// Inc increments the counter by 1.
func (c *Counter) Inc() { atomic.AddUint64(&c.v, 1) }

// Value returns the current counter value.
func (c *Counter) Value() uint64 { return atomic.LoadUint64(&c.v) }

// Metrics holds all named counters for the process lifetime.
type Metrics struct {
	mu       sync.RWMutex
	counters map[string]*Counter
}

// New returns an initialised Metrics registry.
func New() *Metrics {
	return &Metrics{counters: make(map[string]*Counter)}
}

// Inc increments the named counter, creating it if necessary.
func (m *Metrics) Inc(name string) {
	m.mu.Lock()
	c, ok := m.counters[name]
	if !ok {
		c = &Counter{}
		m.counters[name] = c
	}
	m.mu.Unlock()
	c.Inc()
}

// Get returns the current value of the named counter (0 if unknown).
func (m *Metrics) Get(name string) uint64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if c, ok := m.counters[name]; ok {
		return c.Value()
	}
	return 0
}

// Snapshot returns a point-in-time copy of all counters.
func (m *Metrics) Snapshot() map[string]uint64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make(map[string]uint64, len(m.counters))
	for k, c := range m.counters {
		out[k] = c.Value()
	}
	return out
}

// Reset sets every counter back to zero.
func (m *Metrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, c := range m.counters {
		atomic.StoreUint64(&c.v, 0)
	}
}
