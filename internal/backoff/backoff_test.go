package backoff

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("default config invalid: %v", err)
	}
}

func TestValidate_NegativeInitial(t *testing.T) {
	cfg := DefaultConfig()
	cfg.InitialInterval = -1
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative InitialInterval")
	}
}

func TestValidate_MaxLessThanInitial(t *testing.T) {
	cfg := DefaultConfig()
	cfg.MaxInterval = cfg.InitialInterval - 1
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error when MaxInterval < InitialInterval")
	}
}

func TestValidate_MultiplierBelowOne(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Multiplier = 0.5
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for Multiplier < 1")
	}
}

func TestValidate_JitterOutOfRange(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Jitter = 1.5
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for Jitter > 1")
	}
}

func TestNew_InvalidConfig(t *testing.T) {
	cfg := DefaultConfig()
	cfg.InitialInterval = 0
	if _, err := New(cfg); err == nil {
		t.Fatal("expected error from New with invalid config")
	}
}

func TestNext_IncreasesDelay(t *testing.T) {
	cfg := Config{
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     10 * time.Second,
		Multiplier:      2.0,
		Jitter:          0,
	}
	b, err := New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	d0 := b.Next()
	d1 := b.Next()
	if d1 <= d0 {
		t.Errorf("expected d1 (%v) > d0 (%v)", d1, d0)
	}
}

func TestNext_CapsAtMaxInterval(t *testing.T) {
	cfg := Config{
		InitialInterval: 1 * time.Second,
		MaxInterval:     2 * time.Second,
		Multiplier:      10.0,
		Jitter:          0,
	}
	b, _ := New(cfg)
	for i := 0; i < 10; i++ {
		d := b.Next()
		if d > cfg.MaxInterval {
			t.Errorf("attempt %d: delay %v exceeded MaxInterval %v", i, d, cfg.MaxInterval)
		}
	}
}

func TestReset_ResetsAttempt(t *testing.T) {
	cfg := Config{
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     5 * time.Second,
		Multiplier:      2.0,
		Jitter:          0,
	}
	b, _ := New(cfg)
	b.Next()
	b.Next()
	if b.Attempt() != 2 {
		t.Fatalf("expected attempt 2, got %d", b.Attempt())
	}
	b.Reset()
	if b.Attempt() != 0 {
		t.Fatalf("expected attempt 0 after reset, got %d", b.Attempt())
	}
	if d := b.Next(); d != cfg.InitialInterval {
		t.Errorf("after reset first delay should equal InitialInterval, got %v", d)
	}
}
