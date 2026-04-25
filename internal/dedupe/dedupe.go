// Package dedupe provides deduplication of secret maps by tracking seen
// key-value pairs and suppressing repeated entries across fetch cycles.
package dedupe

import "sync"

// Deduplicator tracks previously seen secrets and filters out unchanged entries.
type Deduplicator struct {
	mu   sync.Mutex
	seen map[string]string
}

// New returns a new Deduplicator with an empty seen-set.
func New() *Deduplicator {
	return &Deduplicator{
		seen: make(map[string]string),
	}
}

// Filter returns only the key-value pairs from incoming that differ from
// (or are absent in) the previously seen set. It also updates the seen set
// with all keys present in incoming.
func (d *Deduplicator) Filter(incoming map[string]string) map[string]string {
	d.mu.Lock()
	defer d.mu.Unlock()

	result := make(map[string]string)
	for k, v := range incoming {
		if prev, ok := d.seen[k]; !ok || prev != v {
			result[k] = v
		}
		d.seen[k] = v
	}
	return result
}

// Reset clears the seen set, causing the next Filter call to treat all
// incoming entries as new.
func (d *Deduplicator) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[string]string)
}

// Seen returns a snapshot of the current seen set.
func (d *Deduplicator) Seen() map[string]string {
	d.mu.Lock()
	defer d.mu.Unlock()
	out := make(map[string]string, len(d.seen))
	for k, v := range d.seen {
		out[k] = v
	}
	return out
}
