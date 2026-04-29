package drain_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/drain"
)

func TestNew_DefaultTimeout(t *testing.T) {
	d := drain.New(0)
	if d == nil {
		t.Fatal("expected non-nil Drainer")
	}
}

func TestAcquire_ReturnsFalseAfterClose(t *testing.T) {
	d := drain.New(time.Second)
	d.Wait(context.Background()) // closes immediately — no in-flight work
	if d.Acquire() {
		t.Fatal("expected Acquire to return false after Wait")
	}
}

func TestWait_ReturnsTrueWhenAllRelease(t *testing.T) {
	d := drain.New(2 * time.Second)

	const workers = 5
	var started sync.WaitGroup
	started.Add(workers)

	for i := 0; i < workers; i++ {
		if !d.Acquire() {
			t.Fatalf("Acquire %d failed", i)
		}
		go func() {
			started.Done()
			time.Sleep(20 * time.Millisecond)
			d.Release()
		}()
	}

	started.Wait()
	ctx := context.Background()
	if !d.Wait(ctx) {
		t.Fatal("expected Wait to return true")
	}
}

func TestWait_ReturnsFalseOnContextCancel(t *testing.T) {
	d := drain.New(5 * time.Second)

	if !d.Acquire() {
		t.Fatal("Acquire failed")
	}
	// intentionally never Release — simulate hung operation

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	if d.Wait(ctx) {
		t.Fatal("expected Wait to return false on context cancel")
	}
}

func TestWait_ReturnsFalseOnInternalTimeout(t *testing.T) {
	d := drain.New(30 * time.Millisecond)

	if !d.Acquire() {
		t.Fatal("Acquire failed")
	}
	// intentionally never Release

	if d.Wait(context.Background()) {
		t.Fatal("expected Wait to return false on internal timeout")
	}
}

func TestClosed_FalseBeforeWait(t *testing.T) {
	d := drain.New(time.Second)
	if d.Closed() {
		t.Fatal("expected Closed to be false before Wait")
	}
}

func TestClosed_TrueAfterWait(t *testing.T) {
	d := drain.New(time.Second)
	d.Wait(context.Background())
	if !d.Closed() {
		t.Fatal("expected Closed to be true after Wait")
	}
}
