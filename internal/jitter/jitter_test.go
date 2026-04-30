package jitter_test

import (
	"context"
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/jitter"
)

func TestDefaultConfig(t *testing.T) {
	cfg := jitter.DefaultConfig()
	if cfg.Base <= 0 {
		t.Fatalf("expected positive Base, got %v", cfg.Base)
	}
	if cfg.Factor <= 0 || cfg.Factor > 1 {
		t.Fatalf("expected Factor in (0,1], got %v", cfg.Factor)
	}
}

func TestNew_InvalidConfig(t *testing.T) {
	cases := []struct {
		name string
		cfg  jitter.Config
	}{
		{"zero base", jitter.Config{Base: 0, Factor: 0.1}},
		{"negative base", jitter.Config{Base: -1 * time.Second, Factor: 0.1}},
		{"zero factor", jitter.Config{Base: time.Second, Factor: 0}},
		{"factor above one", jitter.Config{Base: time.Second, Factor: 1.5}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := jitter.New(tc.cfg)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestNext_WithinBounds(t *testing.T) {
	cfg := jitter.Config{Base: 100 * time.Millisecond, Factor: 0.5}
	s, err := jitter.New(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 100; i++ {
		d := s.Next()
		if d < cfg.Base {
			t.Fatalf("duration %v below Base %v", d, cfg.Base)
		}
		max := cfg.Base + time.Duration(float64(cfg.Base)*cfg.Factor)
		if d > max {
			t.Fatalf("duration %v above max %v", d, max)
		}
	}
}

func TestWait_ReturnsOnContextCancel(t *testing.T) {
	cfg := jitter.Config{Base: 10 * time.Second, Factor: 0.1}
	s, err := jitter.New(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	start := time.Now()
	if err := s.Wait(ctx); err == nil {
		t.Fatal("expected context error, got nil")
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("Wait did not return promptly: %v", elapsed)
	}
}

func TestWait_CompletesWithinExpectedWindow(t *testing.T) {
	cfg := jitter.Config{Base: 20 * time.Millisecond, Factor: 0.5}
	s, err := jitter.New(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ctx := context.Background()
	start := time.Now()
	if err := s.Wait(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	elapsed := time.Since(start)
	if elapsed < cfg.Base {
		t.Fatalf("Wait returned too early: %v", elapsed)
	}
}
