package timeout_test

import (
	"context"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/timeout"
)

func TestDefaultConfig(t *testing.T) {
	cfg := timeout.DefaultConfig()
	if cfg.Fetch != 10*time.Second {
		t.Errorf("expected Fetch=10s, got %s", cfg.Fetch)
	}
	if cfg.Exec != 0 {
		t.Errorf("expected Exec=0, got %s", cfg.Exec)
	}
}

func TestValidate_NegativeFetch(t *testing.T) {
	cfg := timeout.Config{Fetch: -1 * time.Second}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative fetch duration")
	}
}

func TestValidate_NegativeExec(t *testing.T) {
	cfg := timeout.Config{Exec: -1 * time.Second, Fetch: 5 * time.Second}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative exec duration")
	}
}

func TestNew_InvalidConfig(t *testing.T) {
	_, err := timeout.New(timeout.Config{Fetch: -1 * time.Second})
	if err == nil {
		t.Fatal("expected error from New with invalid config")
	}
}

func TestExecContext_ZeroTimeout(t *testing.T) {
	e, err := timeout.New(timeout.Config{Fetch: 5 * time.Second, Exec: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ctx, cancel := e.ExecContext(context.Background())
	defer cancel()
	if _, ok := ctx.Deadline(); ok {
		t.Error("expected no deadline when Exec=0")
	}
}

func TestExecContext_WithTimeout(t *testing.T) {
	e, err := timeout.New(timeout.Config{Fetch: 5 * time.Second, Exec: 50 * time.Millisecond})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ctx, cancel := e.ExecContext(context.Background())
	defer cancel()
	if _, ok := ctx.Deadline(); !ok {
		t.Error("expected deadline to be set")
	}
	select {
	case <-ctx.Done():
		if !timeout.IsDeadline(ctx.Err()) {
			t.Errorf("expected deadline error, got %v", ctx.Err())
		}
	case <-time.After(200 * time.Millisecond):
		t.Error("context was not cancelled within expected time")
	}
}

func TestFetchContext_Cancels(t *testing.T) {
	e, err := timeout.New(timeout.Config{Fetch: 30 * time.Millisecond})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ctx, cancel := e.FetchContext(context.Background())
	defer cancel()
	select {
	case <-ctx.Done():
	// ok
	case <-time.After(200 * time.Millisecond):
		t.Error("fetch context did not expire")
	}
}

func TestIsDeadline(t *testing.T) {
	if !timeout.IsDeadline(context.DeadlineExceeded) {
		t.Error("expected IsDeadline=true for context.DeadlineExceeded")
	}
	if timeout.IsDeadline(context.Canceled) {
		t.Error("expected IsDeadline=false for context.Canceled")
	}
}
