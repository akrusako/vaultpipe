// Package checkpoint persists a record of which secret paths have been
// successfully fetched so that vaultpipe can resume after a restart without
// re-reading every path from scratch.
package checkpoint

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Record holds the persisted state for a single secret path.
type Record struct {
	Path      string    `json:"path"`
	FetchedAt time.Time `json:"fetched_at"`
}

// Checkpoint manages a lightweight on-disk JSON file that tracks which
// secret paths have been successfully read.
type Checkpoint struct {
	mu      sync.RWMutex
	path    string
	records map[string]Record
}

// New loads an existing checkpoint file at filePath, or starts with an empty
// state if the file does not yet exist.
func New(filePath string) (*Checkpoint, error) {
	cp := &Checkpoint{
		path:    filePath,
		records: make(map[string]Record),
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return cp, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(data, &cp.records); err != nil {
		return nil, err
	}
	return cp, nil
}

// Mark records that the given secret path was successfully fetched right now.
func (c *Checkpoint) Mark(secretPath string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.records[secretPath] = Record{
		Path:      secretPath,
		FetchedAt: time.Now().UTC(),
	}
	return c.flush()
}

// Has reports whether the given secret path has a checkpoint record.
func (c *Checkpoint) Has(secretPath string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.records[secretPath]
	return ok
}

// Get returns the Record for the given path and whether it was found.
func (c *Checkpoint) Get(secretPath string) (Record, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	r, ok := c.records[secretPath]
	return r, ok
}

// Reset removes all records and overwrites the checkpoint file.
func (c *Checkpoint) Reset() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.records = make(map[string]Record)
	return c.flush()
}

// flush writes the current records map to disk. Caller must hold c.mu.
func (c *Checkpoint) flush() error {
	data, err := json.MarshalIndent(c.records, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.path, data, 0o600)
}
