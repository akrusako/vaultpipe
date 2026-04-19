package signal

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestWatch_CancelsOnSignal(t *testing.T) {
	h := New()

	// Replace notify to immediately send a signal.
	h.notify = func(ch chan<- os.Signal, sigs ...os.Signal) {
		go func() { ch <- sigs[0] }()
	}
	h.stop = func(chan<- os.Signal) {}

	ctx, cancel := h.Watch(context.Background())
	defer cancel()

	select {
	case <-ctx.Done():
		// expected
	case <-time.After(time.Second):
		t.Fatal("context was not cancelled after signal")
	}
}

func TestWatch_CancelsOnParentDone(t *testing.T) {
	h := New()
	h.notify = func(ch chan<- os.Signal, sigs ...os.Signal) {}
	h.stop = func(chan<- os.Signal) {}

	parent, parentCancel := context.WithCancel(context.Background())
	ctx, cancel := h.Watch(parent)
	defer cancel()

	parentCancel()

	select {
	case <-ctx.Done():
		// expected
	case <-time.After(time.Second):
		t.Fatal("context was not cancelled after parent done")
	}
}

func TestNew_DefaultSignals(t *testing.T) {
	h := New()
	if len(h.signals) < 2 {
		t.Fatalf("expected at least 2 signals, got %d", len(h.signals))
	}
}

func TestNew_ExtraSignals(t *testing.T) {
	h := New(os.Interrupt)
	if len(h.signals) != 3 {
		t.Fatalf("expected 3 signals, got %d", len(h.signals))
	}
}
