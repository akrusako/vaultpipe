// Package debounce delays forwarding of secret fetch results until no new
// calls have been made within a configurable quiet period. This prevents
// thundering-herd refreshes when multiple watchers fire in rapid succession.
package debounce

import (
	"context"
	"sync"
	"time"
)

// Fetcher is the function signature used to retrieve secrets.
type Fetcher func(ctx context.Context, paths []string) (map[string]string, error)

// Debouncer wraps a Fetcher and coalesces rapid successive calls.
type Debouncer struct {
	mu      sync.Mutex
	wait    time.Duration
	fetcher Fetcher
	timer   *time.Timer
	pending []string
	result  chan result
}

type result struct {
	data map[string]string
	err  error
}

// New returns a Debouncer that delays execution of fetcher by wait after the
// last call to Fetch. A wait of zero disables debouncing.
func New(wait time.Duration, fetcher Fetcher) (*Debouncer, error) {
	if fetcher == nil {
		return nil, errNilFetcher
	}
	if wait < 0 {
		return nil, errNegativeWait
	}
	return &Debouncer{
		wait:    wait,
		fetcher: fetcher,
		result:  make(chan result, 1),
	}, nil
}

// Fetch schedules a fetch for the given paths. If another call arrives within
// the debounce window the timer resets and only one fetch is performed.
func (d *Debouncer) Fetch(ctx context.Context, paths []string) (map[string]string, error) {
	if d.wait == 0 {
		return d.fetcher(ctx, paths)
	}

	d.mu.Lock()
	d.pending = paths
	if d.timer != nil {
		d.timer.Reset(d.wait)
	} else {
		d.timer = time.AfterFunc(d.wait, func() {
			d.mu.Lock()
			p := d.pending
			d.timer = nil
			d.mu.Unlock()

			data, err := d.fetcher(ctx, p)
			select {
			case d.result <- result{data: data, err: err}:
			default:
			}
		})
	}
	d.mu.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case r := <-d.result:
		return r.data, r.err
	}
}

// errors
var (
	errNilFetcher  = debounceError("debounce: fetcher must not be nil")
	errNegativeWait = debounceError("debounce: wait duration must not be negative")
)

type debounceError string

func (e debounceError) Error() string { return string(e) }
