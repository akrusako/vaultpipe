package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/retry"
)

var errTransient = errors.New("transient error")

func fastConfig() retry.Config {
	return retry.Config{
		MaxAttempts:  3,
		InitialDelay: time.Millisecond,
		MaxDelay:     5 * time.Millisecond,
		Multiplier:   2.0,
	}
}

func TestDo_SuccessOnFirstAttempt(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), fastConfig(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesAndSucceeds(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), fastConfig(), func() error {
		calls++
		if calls < 3 {
			return errTransient
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ExhaustsAttempts(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), fastConfig(), func() error {
		calls++
		return errTransient
	})
	if !errors.Is(err, retry.ErrMaxAttempts) {
		t.Fatalf("expected ErrMaxAttempts, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_PermanentStopsImmediately(t *testing.T) {
	sentinel := errors.New("fatal")
	calls := 0
	err := retry.Do(context.Background(), fastConfig(), func() error {
		calls++
		return retry.Permanent(sentinel)
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := retry.Do(ctx, fastConfig(), func() error {
		return errTransient
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}
