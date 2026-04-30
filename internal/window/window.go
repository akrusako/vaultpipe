// Package window provides a sliding-window counter for rate and frequency
// tracking over a fixed time duration. It is safe for concurrent use.
package window

import (
	"sync"
	"time"
)

// Window tracks event counts within a rolling time window.
type Window struct {
	mu       sync.Mutex
	buckets  []int64
	size     int
	duration time.Duration
	bucket   time.Duration // duration per bucket
	last     time.Time
}

// New creates a Window that spans the given duration split into the given
// number of buckets. Panics if size < 1 or duration <= 0.
func New(duration time.Duration, size int) *Window {
	if size < 1 {
		panic("window: size must be >= 1")
	}
	if duration <= 0 {
		panic("window: duration must be positive")
	}
	return &Window{
		buckets:  make([]int64, size),
		size:     size,
		duration: duration,
		bucket:   duration / time.Duration(size),
		last:     time.Now(),
	}
}

// Add records n events at the current time.
func (w *Window) Add(n int64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.advance(time.Now())
	w.buckets[0] += n
}

// Count returns the total number of events recorded within the window.
func (w *Window) Count() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.advance(time.Now())
	var total int64
	for _, b := range w.buckets {
		total += b
	}
	return total
}

// Reset clears all buckets.
func (w *Window) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	for i := range w.buckets {
		w.buckets[i] = 0
	}
	w.last = time.Now()
}

// advance rotates buckets to account for elapsed time since the last write.
func (w *Window) advance(now time.Time) {
	elapsed := now.Sub(w.last)
	if elapsed < w.bucket {
		return
	}
	shift := int(elapsed / w.bucket)
	if shift >= w.size {
		for i := range w.buckets {
			w.buckets[i] = 0
		}
	} else {
		// rotate right by shift positions
		for i := w.size - 1; i >= shift; i-- {
			w.buckets[i] = w.buckets[i-shift]
		}
		for i := 0; i < shift; i++ {
			w.buckets[i] = 0
		}
	}
	w.last = now
}
