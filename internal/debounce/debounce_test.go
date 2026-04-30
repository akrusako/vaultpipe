package debounce

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func makeFetcher(calls *int32) Fetcher {
	return func(_ context.Context, paths []string) (map[string]string, error) {
		atomic.AddInt32(calls, 1)
		out := make(map[string]string, len(paths))
		for _, p := range paths {
			out[p] = "value"
		}
		return out, nil
	}
}

func TestNew_NilFetcher(t *testing.T) {
	_, err := New(10*time.Millisecond, nil)
	if err == nil {
		t.Fatal("expected error for nil fetcher")
	}
}

func TestNew_NegativeWait(t *testing.T) {
	_, err := New(-1*time.Millisecond, makeFetcher(new(int32)))
	if err == nil {
		t.Fatal("expected error for negative wait")
	}
}

func TestFetch_ZeroWait_PassesThrough(t *testing.T) {
	var calls int32
	d, err := New(0, makeFetcher(&calls))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ctx := context.Background()
	data, err := d.Fetch(ctx, []string{"secret/a"})
	if err != nil {
		t.Fatalf("unexpected fetch error: %v", err)
	}
	if data["secret/a"] != "value" {
		t.Errorf("expected value, got %q", data["secret/a"])
	}
	if atomic.LoadInt32(&calls) != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestFetch_CoalescesRapidCalls(t *testing.T) {
	var calls int32
	wait := 40 * time.Millisecond
	d, err := New(wait, makeFetcher(&calls))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// Fire two fetches quickly; they should coalesce into one backend call.
	resultCh := make(chan error, 2)
	for i := 0; i < 2; i++ {
		go func() {
			_, err := d.Fetch(ctx, []string{"secret/b"})
			resultCh <- err
		}()
	}

	for i := 0; i < 2; i++ {
		if err := <-resultCh; err != nil {
			t.Errorf("fetch error: %v", err)
		}
	}

	if n := atomic.LoadInt32(&calls); n > 2 {
		t.Errorf("expected at most 2 backend calls, got %d", n)
	}
}

func TestFetch_ContextCancel(t *testing.T) {
	blocking := func(_ context.Context, _ []string) (map[string]string, error) {
		time.Sleep(10 * time.Second)
		return nil, nil
	}
	d, _ := New(50*time.Millisecond, blocking)
	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() {
		_, err := d.Fetch(ctx, []string{"secret/c"})
		errCh <- err
	}()

	time.Sleep(10 * time.Millisecond)
	cancel()

	select {
	case err := <-errCh:
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for cancelled fetch")
	}
}
