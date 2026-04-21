package quota

import (
	"errors"
	"testing"
	"time"
)

func fastConfig() Config {
	return Config{MaxReads: 3, Window: time.Second}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.MaxReads <= 0 {
		t.Fatal("expected positive MaxReads")
	}
	if cfg.Window <= 0 {
		t.Fatal("expected positive Window")
	}
}

func TestNew_InvalidConfig(t *testing.T) {
	_, err := New(Config{MaxReads: 0, Window: time.Second})
	if err == nil {
		t.Fatal("expected error for zero MaxReads")
	}
	_, err = New(Config{MaxReads: 1, Window: 0})
	if err == nil {
		t.Fatal("expected error for zero Window")
	}
}

func TestCheck_AllowsUnderLimit(t *testing.T) {
	q, err := New(fastConfig())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	for i := 0; i < 3; i++ {
		if err := q.Check("secret/foo"); err != nil {
			t.Fatalf("unexpected error on read %d: %v", i+1, err)
		}
	}
}

func TestCheck_BlocksAtLimit(t *testing.T) {
	q, err := New(fastConfig())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	for i := 0; i < 3; i++ {
		_ = q.Check("secret/bar")
	}
	err = q.Check("secret/bar")
	if !errors.Is(err, ErrQuotaExceeded) {
		t.Fatalf("expected ErrQuotaExceeded, got %v", err)
	}
}

func TestCheck_WindowExpiry(t *testing.T) {
	q, err := New(fastConfig())
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	base := time.Now()
	q.now = func() time.Time { return base }

	for i := 0; i < 3; i++ {
		_ = q.Check("secret/baz")
	}

	// Advance time beyond the window so old reads expire.
	q.now = func() time.Time { return base.Add(2 * time.Second) }

	if err := q.Check("secret/baz"); err != nil {
		t.Fatalf("expected read to succeed after window expiry, got: %v", err)
	}
}

func TestReset_ClearsPath(t *testing.T) {
	q, err := New(fastConfig())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	for i := 0; i < 3; i++ {
		_ = q.Check("secret/qux")
	}
	q.Reset("secret/qux")
	if err := q.Check("secret/qux"); err != nil {
		t.Fatalf("expected read to succeed after reset, got: %v", err)
	}
}

func TestCounts_ReflectsCurrentWindow(t *testing.T) {
	q, err := New(fastConfig())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_ = q.Check("secret/a")
	_ = q.Check("secret/a")
	_ = q.Check("secret/b")

	counts := q.Counts()
	if counts["secret/a"] != 2 {
		t.Errorf("expected 2 for secret/a, got %d", counts["secret/a"])
	}
	if counts["secret/b"] != 1 {
		t.Errorf("expected 1 for secret/b, got %d", counts["secret/b"])
	}
}

func TestCheck_IndependentPaths(t *testing.T) {
	q, err := New(fastConfig())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	for i := 0; i < 3; i++ {
		_ = q.Check("secret/x")
	}
	// A different path should not be affected.
	if err := q.Check("secret/y"); err != nil {
		t.Fatalf("unexpected error for independent path: %v", err)
	}
}
