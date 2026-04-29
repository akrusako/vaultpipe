// Package drain provides a graceful shutdown helper that waits for
// in-flight secret fetch and exec operations to complete before exit.
package drain

import (
	"context"
	"sync"
	"time"
)

// Drainer tracks active operations and blocks until all finish or the
// deadline is exceeded.
type Drainer struct {
	mu      sync.Mutex
	wg      sync.WaitGroup
	closed  bool
	timeout time.Duration
}

// New returns a Drainer with the given idle timeout.
func New(timeout time.Duration) *Drainer {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &Drainer{timeout: timeout}
}

// Acquire registers a new in-flight operation.
// It returns false if the Drainer has already been closed.
func (d *Drainer) Acquire() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.closed {
		return false
	}
	d.wg.Add(1)
	return true
}

// Release marks one in-flight operation as complete.
func (d *Drainer) Release() {
	d.wg.Done()
}

// Wait closes the Drainer to new acquisitions and blocks until all
// in-flight operations finish or ctx is cancelled / the internal
// timeout expires, whichever comes first.
// It returns true if all operations finished cleanly.
func (d *Drainer) Wait(ctx context.Context) bool {
	d.mu.Lock()
	d.closed = true
	d.mu.Unlock()

	done := make(chan struct{})
	go func() {
		d.wg.Wait()
		close(done)
	}()

	timer := time.NewTimer(d.timeout)
	defer timer.Stop()

	select {
	case <-done:
		return true
	case <-ctx.Done():
		return false
	case <-timer.C:
		return false
	}
}

// Closed reports whether the Drainer has been closed to new work.
func (d *Drainer) Closed() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.closed
}
