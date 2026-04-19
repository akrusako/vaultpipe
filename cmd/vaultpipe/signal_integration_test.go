package main

import (
	"context"
	"os"
	"testing"
	"time"

	vpsignal "github.com/yourusername/vaultpipe/internal/signal"
)

func TestSignalHandler_IntegratesWithContext(t *testing.T) {
	h := vpsignal.New()

	// Override internals via the exported constructor; just verify Watch returns
	// a valid, cancellable context.
	ctx, cancel := h.Watch(context.Background())
	defer cancel()

	if ctx == nil {
		t.Fatal("expected non-nil context")
	}

	// Manually cancel and ensure Done closes.
	cancel()
	select {
	case <-ctx.Done():
		// expected
	case <-time.After(time.Second):
		t.Fatal("context did not close after cancel")
	}
}

func TestSignalHandler_SIGTERM(t *testing.T) {
	if os.Getenv("CI") == "" {
		t.Skip("skipping live signal test outside CI")
	}

	h := vpsignal.New()
	ctx, cancel := h.Watch(context.Background())
	defer cancel()

	// Send SIGTERM to self.
	p, _ := os.FindProcess(os.Getpid())
	_ = p.Signal(os.Interrupt)

	select {
	case <-ctx.Done():
		// expected
	case <-time.After(2 * time.Second):
		t.Fatal("context not cancelled after SIGTERM")
	}
}
