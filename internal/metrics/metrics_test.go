package metrics_test

import (
	"sync"
	"testing"

	"github.com/yourusername/vaultpipe/internal/metrics"
)

func TestGet_UnknownCounterIsZero(t *testing.T) {
	m := metrics.New()
	if v := m.Get("missing"); v != 0 {
		t.Fatalf("expected 0, got %d", v)
	}
}

func TestInc_CreatesAndIncrementsCounter(t *testing.T) {
	m := metrics.New()
	m.Inc("vault.fetch")
	m.Inc("vault.fetch")
	if v := m.Get("vault.fetch"); v != 2 {
		t.Fatalf("expected 2, got %d", v)
	}
}

func TestSnapshot_ContainsAllCounters(t *testing.T) {
	m := metrics.New()
	m.Inc("a")
	m.Inc("b")
	m.Inc("b")
	snap := m.Snapshot()
	if snap["a"] != 1 {
		t.Errorf("a: expected 1, got %d", snap["a"])
	}
	if snap["b"] != 2 {
		t.Errorf("b: expected 2, got %d", snap["b"])
	}
}

func TestSnapshot_IsCopy(t *testing.T) {
	m := metrics.New()
	m.Inc("x")
	snap := m.Snapshot()
	m.Inc("x") // mutate after snapshot
	if snap["x"] != 1 {
		t.Errorf("snapshot should not reflect later increments")
	}
}

func TestReset_ZeroesCounters(t *testing.T) {
	m := metrics.New()
	m.Inc("exec.start")
	m.Inc("exec.start")
	m.Reset()
	if v := m.Get("exec.start"); v != 0 {
		t.Fatalf("expected 0 after reset, got %d", v)
	}
}

func TestInc_ConcurrentSafe(t *testing.T) {
	m := metrics.New()
	const goroutines = 50
	const incsEach = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < incsEach; j++ {
				m.Inc("concurrent")
			}
		}()
	}
	wg.Wait()
	want := uint64(goroutines * incsEach)
	if v := m.Get("concurrent"); v != want {
		t.Fatalf("expected %d, got %d", want, v)
	}
}
