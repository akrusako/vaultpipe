package circuit_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/circuit"
)

func TestDefaultConfig(t *testing.T) {
	cfg := circuit.DefaultConfig()
	if cfg.MaxFailures <= 0 {
		t.Errorf("expected MaxFailures > 0, got %d", cfg.MaxFailures)
	}
	if cfg.ResetTimeout <= 0 {
		t.Errorf("expected ResetTimeout > 0, got %v", cfg.ResetTimeout)
	}
}

func TestNew_InvalidConfig(t *testing.T) {
	_, err := circuit.New(circuit.Config{MaxFailures: 0, ResetTimeout: time.Second})
	if err == nil {
		t.Fatal("expected error for zero MaxFailures")
	}
	_, err = circuit.New(circuit.Config{MaxFailures: 3, ResetTimeout: 0})
	if err == nil {
		t.Fatal("expected error for zero ResetTimeout")
	}
}

func TestBreaker_InitiallyClosed(t *testing.T) {
	b, err := circuit.New(circuit.DefaultConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.State() != circuit.StateClosed {
		t.Errorf("expected StateClosed, got %v", b.State())
	}
	if !b.Allow() {
		t.Error("expected Allow() == true when closed")
	}
}

func TestBreaker_OpensAfterMaxFailures(t *testing.T) {
	cfg := circuit.Config{MaxFailures: 3, ResetTimeout: time.Minute}
	b, _ := circuit.New(cfg)

	for i := 0; i < 3; i++ {
		if b.State() == circuit.StateOpen {
			t.Fatalf("opened too early at failure %d", i)
		}
		b.RecordFailure()
	}
	if b.State() != circuit.StateOpen {
		t.Errorf("expected StateOpen after max failures, got %v", b.State())
	}
	if b.Allow() {
		t.Error("expected Allow() == false when open")
	}
}

func TestBreaker_HalfOpenAfterTimeout(t *testing.T) {
	cfg := circuit.Config{MaxFailures: 1, ResetTimeout: 10 * time.Millisecond}
	b, _ := circuit.New(cfg)
	b.RecordFailure()

	time.Sleep(20 * time.Millisecond)

	if !b.Allow() {
		t.Error("expected Allow() == true after reset timeout (half-open)")
	}
	if b.State() != circuit.StateHalfOpen {
		t.Errorf("expected StateHalfOpen, got %v", b.State())
	}
}

func TestBreaker_RecordSuccess_ResetsToClosed(t *testing.T) {
	cfg := circuit.Config{MaxFailures: 1, ResetTimeout: 10 * time.Millisecond}
	b, _ := circuit.New(cfg)
	b.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	b.Allow() // transitions to half-open
	b.RecordSuccess()

	if b.State() != circuit.StateClosed {
		t.Errorf("expected StateClosed after success, got %v", b.State())
	}
}
