package lease

import (
	"context"
	"time"
)

// WarnFunc is called when a lease is approaching expiry.
type WarnFunc func(info Info)

// WatcherConfig configures the expiry watcher.
type WatcherConfig struct {
	// Interval controls how often the tracker is polled.
	Interval time.Duration
	// Threshold is the TTL below which a lease is considered expiring.
	Threshold time.Duration
	// OnExpiring is called for each lease below the threshold.
	OnExpiring WarnFunc
}

// DefaultWatcherConfig returns sensible defaults.
func DefaultWatcherConfig() WatcherConfig {
	return WatcherConfig{
		Interval:  30 * time.Second,
		Threshold: 5 * time.Minute,
	}
}

// Watcher polls a Tracker and invokes a callback for expiring leases.
type Watcher struct {
	cfg     WatcherConfig
	tracker *Tracker
}

// NewWatcher creates a Watcher backed by the given Tracker.
func NewWatcher(tracker *Tracker, cfg WatcherConfig) *Watcher {
	if cfg.Interval <= 0 {
		cfg.Interval = DefaultWatcherConfig().Interval
	}
	if cfg.Threshold <= 0 {
		cfg.Threshold = DefaultWatcherConfig().Threshold
	}
	return &Watcher{cfg: cfg, tracker: tracker}
}

// Watch starts the polling loop. It blocks until ctx is cancelled.
func (w *Watcher) Watch(ctx context.Context) {
	ticker := time.NewTicker(w.cfg.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.check()
		}
	}
}

func (w *Watcher) check() {
	if w.cfg.OnExpiring == nil {
		return
	}
	for _, info := range w.tracker.Expiring(w.cfg.Threshold) {
		w.cfg.OnExpiring(info)
	}
}
