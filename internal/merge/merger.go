// Package merge provides utilities for combining multiple secret maps
// into a single environment map with configurable conflict resolution.
package merge

import "fmt"

// Strategy controls how key conflicts are resolved when merging maps.
type Strategy int

const (
	// StrategyFirst keeps the value from the first map that defines a key.
	StrategyFirst Strategy = iota
	// StrategyLast keeps the value from the last map that defines a key.
	StrategyLast
	// StrategyError returns an error when a duplicate key is encountered.
	StrategyError
)

// Merger combines multiple secret maps into one.
type Merger struct {
	strategy Strategy
}

// New returns a Merger using the given conflict resolution strategy.
func New(s Strategy) *Merger {
	return &Merger{strategy: s}
}

// Merge combines the provided maps in order. Conflict resolution follows
// the strategy set on the Merger. The returned map is always a new
// allocation and does not alias any input.
func (m *Merger) Merge(maps ...map[string]string) (map[string]string, error) {
	out := make(map[string]string)

	for _, src := range maps {
		for k, v := range src {
			if existing, ok := out[k]; ok {
				switch m.strategy {
				case StrategyFirst:
					// keep existing — do nothing
					_ = existing
				case StrategyLast:
					out[k] = v
				case StrategyError:
					return nil, fmt.Errorf("merge: duplicate key %q", k)
				}
			} else {
				out[k] = v
			}
		}
	}

	return out, nil
}

// MustMerge is like Merge but panics on error. Useful in tests.
func (m *Merger) MustMerge(maps ...map[string]string) map[string]string {
	out, err := m.Merge(maps...)
	if err != nil {
		panic(err)
	}
	return out
}
