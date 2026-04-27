// Package coalesce provides a secret source that returns the first non-empty
// value for a given key across an ordered list of secret maps.
package coalesce

import "errors"

// ErrNoValue is returned when no source contains a non-empty value for a key.
var ErrNoValue = errors.New("coalesce: no non-empty value found for key")

// Coalescer selects the first non-empty value for each key from an ordered
// list of secret maps. Sources are evaluated left-to-right; the first map
// that contains a non-empty string for the requested key wins.
type Coalescer struct {
	sources []map[string]string
}

// New returns a Coalescer that will search the given sources in order.
// Sources are stored by reference; callers must not modify them after New.
func New(sources ...map[string]string) *Coalescer {
	copied := make([]map[string]string, len(sources))
	copy(copied, sources)
	return &Coalescer{sources: copied}
}

// Lookup returns the first non-empty value for key across all sources.
// It returns ErrNoValue if every source either lacks the key or holds an
// empty string.
func (c *Coalescer) Lookup(key string) (string, error) {
	for _, src := range c.sources {
		if v, ok := src[key]; ok && v != "" {
			return v, nil
		}
	}
	return "", ErrNoValue
}

// Merge builds a single flat map by applying Lookup to the union of all keys
// found across every source. Keys whose coalesced value is empty are omitted.
func (c *Coalescer) Merge() map[string]string {
	seen := make(map[string]struct{})
	for _, src := range c.sources {
		for k := range src {
			seen[k] = struct{}{}
		}
	}
	out := make(map[string]string, len(seen))
	for k := range seen {
		if v, err := c.Lookup(k); err == nil {
			out[k] = v
		}
	}
	return out
}

// Add appends an additional source at the lowest priority (end of the list).
func (c *Coalescer) Add(src map[string]string) {
	c.sources = append(c.sources, src)
}
