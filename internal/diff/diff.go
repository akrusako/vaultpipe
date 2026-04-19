// Package diff compares two sets of secrets and reports which keys were
// added, removed, or changed between rotations.
package diff

// Change describes a single key-level difference between two secret maps.
type Change struct {
	Key    string
	Action string // "added", "removed", "changed"
}

// Differ compares secret snapshots.
type Differ struct{}

// New returns a new Differ.
func New() *Differ {
	return &Differ{}
}

// Compare returns the list of changes between prev and next secret maps.
// Values are compared but never stored in Change to avoid leaking secrets.
func (d *Differ) Compare(prev, next map[string]string) []Change {
	var changes []Change

	for k, nv := range next {
		pv, ok := prev[k]
		if !ok {
			changes = append(changes, Change{Key: k, Action: "added"})
		} else if pv != nv {
			changes = append(changes, Change{Key: k, Action: "changed"})
		}
	}

	for k := range prev {
		if _, ok := next[k]; !ok {
			changes = append(changes, Change{Key: k, Action: "removed"})
		}
	}

	return changes
}

// HasChanges returns true when prev and next differ in any key or value.
func (d *Differ) HasChanges(prev, next map[string]string) bool {
	return len(d.Compare(prev, next)) > 0
}
