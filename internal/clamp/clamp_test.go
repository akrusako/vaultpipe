package clamp_test

import (
	"context"
	"errors"
	"testing"

	"github.com/yourusername/vaultpipe/internal/clamp"
)

func TestNew_InvalidConfig(t *testing.T) {
	_, err := clamp.New(clamp.Options{MinBytes: 10, MaxBytes: 5})
	if err == nil {
		t.Fatal("expected error when MaxBytes < MinBytes")
	}
}

func TestNew_ValidConfig(t *testing.T) {
	_, err := clamp.New(clamp.Options{MinBytes: 5, MaxBytes: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApply_PassesThroughWithinBounds(t *testing.T) {
	c, _ := clamp.New(clamp.Options{MinBytes: 2, MaxBytes: 8})
	secrets := map[string]string{"KEY": "hello"}
	out, err := c.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "hello" {
		t.Fatalf("expected hello, got %q", out["KEY"])
	}
}

func TestApply_ErrorOnShortValue(t *testing.T) {
	c, _ := clamp.New(clamp.Options{MinBytes: 10})
	_, err := c.Apply(map[string]string{"K": "hi"})
	if err == nil {
		t.Fatal("expected error for short value")
	}
}

func TestApply_DropShort(t *testing.T) {
	c, _ := clamp.New(clamp.Options{MinBytes: 10, DropShort: true})
	out, err := c.Apply(map[string]string{"K": "hi", "LONG": "longenoughvalue"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["K"]; ok {
		t.Error("expected short key to be dropped")
	}
	if out["LONG"] != "longenoughvalue" {
		t.Error("expected long-enough key to be kept")
	}
}

func TestApply_ErrorOnLongValue(t *testing.T) {
	c, _ := clamp.New(clamp.Options{MaxBytes: 3})
	_, err := c.Apply(map[string]string{"K": "toolong"})
	if err == nil {
		t.Fatal("expected error for long value")
	}
}

func TestApply_DropLong(t *testing.T) {
	c, _ := clamp.New(clamp.Options{MaxBytes: 3, DropLong: true})
	out, err := c.Apply(map[string]string{"K": "toolong", "S": "ok"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["K"]; ok {
		t.Error("expected long key to be dropped")
	}
	if out["S"] != "ok" {
		t.Error("expected short key to be kept")
	}
}

func TestNewMiddleware_PropagatesFetchError(t *testing.T) {
	c, _ := clamp.New(clamp.Options{})
	fetchErr := errors.New("vault unavailable")
	fetcher := clamp.NewMiddleware(c, func(_ context.Context, _ []string) (map[string]string, error) {
		return nil, fetchErr
	})
	_, err := fetcher(context.Background(), []string{"secret/data/app"})
	if !errors.Is(err, fetchErr) {
		t.Fatalf("expected fetch error, got %v", err)
	}
}

func TestNewMiddleware_AppliesConstraints(t *testing.T) {
	c, _ := clamp.New(clamp.Options{MaxBytes: 4, DropLong: true})
	fetcher := clamp.NewMiddleware(c, func(_ context.Context, _ []string) (map[string]string, error) {
		return map[string]string{"A": "hi", "B": "waytoolong"}, nil
	})
	out, err := fetcher(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["B"]; ok {
		t.Error("expected B to be dropped by middleware")
	}
	if out["A"] != "hi" {
		t.Error("expected A to pass through middleware")
	}
}
