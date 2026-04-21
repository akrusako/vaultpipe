package ratelimit_test

import (
	"context"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/ratelimit"
)

func TestDefaultConfig(t *testing.T) {
	cfg := ratelimit.DefaultConfig()
	if cfg.Rate <= 0 {
		t.Fatalf("expected positive rate, got %v", cfg.Rate)
	}
	if cfg.Burst <= 0 {
		t.Fatalf("expected positive burst, got %v", cfg.Burst)
	}
}

func TestValidate_ZeroRate(t *testing.T) {
	cfg := ratelimit.Config{Rate: 0, Burst: 5}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero rate")
	}
}

func TestValidate_ZeroBurst(t *testing.T) {
	cfg := ratelimit.Config{Rate: 1, Burst: 0}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero burst")
	}
}

func TestNew_InvalidConfig(t *testing.T) {
	_, err := ratelimit.New(ratelimit.Config{Rate: -1, Burst: 5})
	if err == nil {
		t.Fatal("expected error for invalid config")
	}
}

func TestWait_ConsumesToken(t *testing.T) {
	lim, err := ratelimit.New(ratelimit.Config{Rate: 100, Burst: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ctx := context.Background()
	if err := lim.Wait(ctx); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestWait_RespectsContextCancellation(t *testing.T) {
	// Very low rate so tokens exhaust quickly.
	lim, err := ratelimit.New(ratelimit.Config{Rate: 0.001, Burst: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Drain the single burst token.
	ctx := context.Background()
	_ = lim.Wait(ctx)

	// Now cancel immediately — Wait should return ctx.Err().
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	start := time.Now()
	err = lim.Wait(ctx)
	if err == nil {
		t.Fatal("expected context error")
	}
	if time.Since(start) > 500*time.Millisecond {
		t.Fatal("Wait blocked too long after context cancellation")
	}
}

func TestWait_BurstAllowsMultiple(t *testing.T) {
	lim, err := ratelimit.New(ratelimit.Config{Rate: 100, Burst: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ctx := context.Background()
	for i := 0; i < 5; i++ {
		if err := lim.Wait(ctx); err != nil {
			t.Fatalf("call %d: unexpected error: %v", i, err)
		}
	}
}
