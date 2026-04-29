package fence_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/fence"
)

func TestOpen_Success(t *testing.T) {
	f := fence.New()
	if err := f.Open("tok1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.IsOpen() {
		t.Fatal("expected fence to be open")
	}
}

func TestOpen_AlreadyOpen(t *testing.T) {
	f := fence.New()
	_ = f.Open("tok1")
	if err := f.Open("tok2"); err != fence.ErrAlreadyOpen {
		t.Fatalf("expected ErrAlreadyOpen, got %v", err)
	}
}

func TestClose_NotOpen(t *testing.T) {
	f := fence.New()
	if err := f.Close("tok1"); err != fence.ErrNotOpen {
		t.Fatalf("expected ErrNotOpen, got %v", err)
	}
}

func TestClose_TokenMismatch(t *testing.T) {
	f := fence.New()
	_ = f.Open("tok1")
	if err := f.Close("wrong"); err == nil {
		t.Fatal("expected token mismatch error")
	}
	// fence should still be open after mismatch
	if !f.IsOpen() {
		t.Fatal("fence should remain open after token mismatch")
	}
}

func TestClose_Success(t *testing.T) {
	f := fence.New()
	_ = f.Open("tok1")
	if err := f.Close("tok1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.IsOpen() {
		t.Fatal("expected fence to be closed")
	}
}

func TestClose_Expired(t *testing.T) {
	now := time.Now()
	f := fence.New(fence.WithTTL(10 * time.Millisecond))

	// Patch internal clock via a sub-test helper approach: use real sleep.
	_ = f.Open("tok1")
	time.Sleep(20 * time.Millisecond)

	if err := f.Close("tok1"); err != fence.ErrExpired {
		t.Fatalf("expected ErrExpired, got %v (now=%v)", err, now)
	}
	if f.IsOpen() {
		t.Fatal("fence should be closed after expiry")
	}
}

func TestReopenAfterClose(t *testing.T) {
	f := fence.New()
	_ = f.Open("tok1")
	_ = f.Close("tok1")
	if err := f.Open("tok2"); err != nil {
		t.Fatalf("expected fence to be reopenable, got %v", err)
	}
}
