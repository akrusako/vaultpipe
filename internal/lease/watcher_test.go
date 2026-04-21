package lease

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestWatcher_CallsOnExpiring(t *testing.T) {
	tr := New()
	// Add a lease that is already almost expired.
	expiring := Info{
		LeaseID:  "exp-1",
		Path:     "secret/exp",
		Duration: 2 * time.Second,
		Acquired: time.Now().Add(-1900 * time.Millisecond),
	}
	tr.Add(expiring)

	var mu sync.Mutex
	var called []string

	cfg := WatcherConfig{
		Interval:  20 * time.Millisecond,
		Threshold: 5 * time.Second,
		OnExpiring: func(info Info) {
			mu.Lock()
			called = append(called, info.LeaseID)
			mu.Unlock()
		},
	}

	w := NewWatcher(tr, cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	w.Watch(ctx)

	mu.Lock()
	defer mu.Unlock()
	if len(called) == 0 {
		t.Fatal("expected OnExpiring to be called at least once")
	}
	if called[0] != "exp-1" {
		t.Errorf("wrong lease ID: %q", called[0])
	}
}

func TestWatcher_NoCallbackWhenNilHandler(t *testing.T) {
	tr := New()
	tr.Add(Info{
		LeaseID:  "x",
		Path:     "secret/x",
		Duration: time.Second,
		Acquired: time.Now().Add(-999 * time.Millisecond),
	})
	cfg := WatcherConfig{
		Interval:   20 * time.Millisecond,
		Threshold:  5 * time.Second,
		OnExpiring: nil,
	}
	w := NewWatcher(tr, cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Millisecond)
	defer cancel()
	// Should not panic.
	w.Watch(ctx)
}

func TestWatcher_CancelsCleanly(t *testing.T) {
	tr := New()
	cfg := DefaultWatcherConfig()
	cfg.Interval = 10 * time.Millisecond
	w := NewWatcher(tr, cfg)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		w.Watch(ctx)
		close(done)
	}()
	cancel()
	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Watch did not return after context cancellation")
	}
}

func TestDefaultWatcherConfig(t *testing.T) {
	cfg := DefaultWatcherConfig()
	if cfg.Interval != 30*time.Second {
		t.Errorf("Interval: got %v want 30s", cfg.Interval)
	}
	if cfg.Threshold != 5*time.Minute {
		t.Errorf("Threshold: got %v want 5m", cfg.Threshold)
	}
}
