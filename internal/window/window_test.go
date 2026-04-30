package window_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/window"
)

func TestNew_PanicsOnZeroSize(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for size=0")
		}
	}()
	window.New(time.Second, 0)
}

func TestNew_PanicsOnNegativeDuration(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for negative duration")
		}
	}()
	window.New(-time.Second, 4)
}

func TestCount_StartsAtZero(t *testing.T) {
	w := window.New(time.Second, 4)
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestAdd_AccumulatesCount(t *testing.T) {
	w := window.New(time.Second, 4)
	w.Add(3)
	w.Add(7)
	if got := w.Count(); got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
}

func TestReset_ZeroesCount(t *testing.T) {
	w := window.New(time.Second, 4)
	w.Add(5)
	w.Reset()
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestCount_DropsOldBuckets(t *testing.T) {
	// Use a very short window so we can observe expiry without long sleeps.
	dur := 40 * time.Millisecond
	w := window.New(dur, 4)
	w.Add(100)
	// Wait for the entire window to elapse.
	time.Sleep(dur + 10*time.Millisecond)
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0 after window expired, got %d", got)
	}
}

func TestAdd_IsThreadSafe(t *testing.T) {
	w := window.New(time.Second, 4)
	done := make(chan struct{})
	for i := 0; i < 50; i++ {
		go func() {
			w.Add(1)
			done <- struct{}{}
		}()
	}
	for i := 0; i < 50; i++ {
		<-done
	}
	if got := w.Count(); got != 50 {
		t.Fatalf("expected 50, got %d", got)
	}
}
