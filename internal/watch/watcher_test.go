package watch_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/watch"
)

type mockFetcher struct {
	mu      sync.Mutex
	calls   int
	results []map[string]string
	err     error
}

func (m *mockFetcher) ReadSecrets(_ context.Context, _ []string) (map[string]string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.err != nil {
		return nil, m.err
	}
	idx := m.calls
	if idx >= len(m.results) {
		idx = len(m.results) - 1
	}
	m.calls++
	return m.results[idx], nil
}

func TestWatch_DetectsChange(t *testing.T) {
	fetcher := &mockFetcher{
		results: []map[string]string{
			{"KEY": "v1"},
			{"KEY": "v2"},
		},
	}

	var mu sync.Mutex
	var gotAdded, gotChanged map[string]string

	w := watch.New(fetcher, []string{"secret/app"}, 20*time.Millisecond, func(added, removed, changed map[string]string) {
		mu.Lock()
		gotAdded = added
		gotChanged = changed
		mu.Unlock()
	})

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	w.Watch(ctx) //nolint:errcheck

	mu.Lock()
	defer mu.Unlock()
	_ = gotAdded
	if gotChanged["KEY"] != "v2" {
		t.Errorf("expected changed KEY=v2, got %q", gotChanged["KEY"])
	}
}

func TestWatch_NoCallbackOnNoChange(t *testing.T) {
	fetcher := &mockFetcher{
		results: []map[string]string{{"KEY": "stable"}},
	}

	called := 0
	w := watch.New(fetcher, []string{"secret/app"}, 20*time.Millisecond, func(_, _, _ map[string]string) {
		called++
	})

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()
	w.Watch(ctx) //nolint:errcheck

	if called > 1 {
		t.Errorf("onChange called %d times on identical secrets, want <=1", called)
	}
}

func TestWatch_ToleratesFetchError(t *testing.T) {
	fetcher := &mockFetcher{err: errors.New("vault unavailable")}

	w := watch.New(fetcher, []string{"secret/app"}, 20*time.Millisecond, func(_, _, _ map[string]string) {
		t.Error("onChange should not be called on fetch error")
	})

	ctx, cancel := context.WithTimeout(context.Background(), 70*time.Millisecond)
	defer cancel()
	w.Watch(ctx) //nolint:errcheck
}

func TestWatch_CancelsCleanly(t *testing.T) {
	fetcher := &mockFetcher{results: []map[string]string{{"K": "v"}}}
	w := watch.New(fetcher, []string{"secret/app"}, 10*time.Millisecond, nil)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- w.Watch(ctx) }()
	cancel()

	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Watch did not return after context cancel")
	}
}
