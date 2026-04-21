package throttle_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/throttle"
)

func TestDefaultConfig(t *testing.T) {
	cfg := throttle.DefaultConfig()
	if cfg.Concurrency != 8 {
		t.Fatalf("expected concurrency 8, got %d", cfg.Concurrency)
	}
}

func TestValidate_ZeroConcurrency(t *testing.T) {
	cfg := throttle.Config{Concurrency: 0}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero concurrency")
	}
}

func TestNew_InvalidConfig(t *testing.T) {
	_, err := throttle.New(throttle.Config{Concurrency: -1})
	if err == nil {
		t.Fatal("expected error for negative concurrency")
	}
}

func TestAcquireRelease_Basic(t *testing.T) {
	l, err := throttle.New(throttle.Config{Concurrency: 2})
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	if err := l.Acquire(ctx); err != nil {
		t.Fatal(err)
	}
	if l.Inflight() != 1 {
		t.Fatalf("expected 1 inflight, got %d", l.Inflight())
	}
	l.Release()
	if l.Inflight() != 0 {
		t.Fatalf("expected 0 inflight after release, got %d", l.Inflight())
	}
}

func TestAcquire_BlocksAtLimit(t *testing.T) {
	l, _ := throttle.New(throttle.Config{Concurrency: 1})
	ctx := context.Background()
	if err := l.Acquire(ctx); err != nil {
		t.Fatal(err)
	}

	ctx2, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := l.Acquire(ctx2)
	if err == nil {
		t.Fatal("expected throttle error when limit reached")
	}
}

func TestAcquire_ConcurrentSafe(t *testing.T) {
	const workers = 20
	const limit = 5
	l, _ := throttle.New(throttle.Config{Concurrency: limit})

	var mu sync.Mutex
	var peak int
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = l.Acquire(context.Background())
			mu.Lock()
			if l.Inflight() > peak {
				peak = l.Inflight()
			}
			mu.Unlock()
			time.Sleep(5 * time.Millisecond)
			l.Release()
		}()
	}
	wg.Wait()
	if peak > limit {
		t.Fatalf("peak inflight %d exceeded limit %d", peak, limit)
	}
}
