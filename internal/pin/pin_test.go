package pin_test

import (
	"context"
	"errors"
	"testing"

	"github.com/yourusername/vaultpipe/internal/pin"
)

type mockFetcher struct {
	result map[string]string
	err    error
}

func (m *mockFetcher) Fetch(_ context.Context, _ []string) (map[string]string, error) {
	return m.result, m.err
}

func TestNew_NilFetcherReturnsError(t *testing.T) {
	_, err := pin.New(nil, pin.Config{})
	if err == nil {
		t.Fatal("expected error for nil fetcher")
	}
}

func TestNew_NegativeMaxDriftReturnsError(t *testing.T) {
	f := &mockFetcher{result: map[string]string{}}
	_, err := pin.New(f, pin.Config{MaxDrift: -1})
	if err == nil {
		t.Fatal("expected error for negative MaxDrift")
	}
}

func TestFetch_FirstCallEstablishesBaseline(t *testing.T) {
	f := &mockFetcher{result: map[string]string{"KEY": "val"}}
	p, _ := pin.New(f, pin.Config{MaxDrift: 0})

	got, err := p.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY"] != "val" {
		t.Fatalf("expected val, got %q", got["KEY"])
	}
}

func TestFetch_NoDriftAllowed_SameMapPasses(t *testing.T) {
	f := &mockFetcher{result: map[string]string{"A": "1", "B": "2"}}
	p, _ := pin.New(f, pin.Config{MaxDrift: 0})

	p.Fetch(context.Background(), nil) // establish baseline
	_, err := p.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("expected no error for identical map, got: %v", err)
	}
}

func TestFetch_ExceedsDrift_ReturnsViolation(t *testing.T) {
	calls := 0
	f := &mockFetcher{}
	p, _ := pin.New(f, pin.Config{MaxDrift: 1})

	f.result = map[string]string{"A": "1", "B": "2"}
	p.Fetch(context.Background(), nil) // baseline
	calls++

	f.result = map[string]string{"A": "changed", "B": "also-changed"}
	_, err := p.Fetch(context.Background(), nil)
	if !errors.Is(err, pin.ErrPinViolation) {
		t.Fatalf("expected ErrPinViolation, got: %v", err)
	}
	_ = calls
}

func TestFetch_WithinDriftLimit_Passes(t *testing.T) {
	f := &mockFetcher{result: map[string]string{"A": "1", "B": "2"}}
	p, _ := pin.New(f, pin.Config{MaxDrift: 1})

	p.Fetch(context.Background(), nil) // baseline

	f.result = map[string]string{"A": "changed", "B": "2"} // 1 key changed
	_, err := p.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("expected no error within drift limit, got: %v", err)
	}
}

func TestFetch_InnerError_Propagates(t *testing.T) {
	f := &mockFetcher{err: errors.New("vault unavailable")}
	p, _ := pin.New(f, pin.Config{})

	_, err := p.Fetch(context.Background(), nil)
	if err == nil || err.Error() != "vault unavailable" {
		t.Fatalf("expected inner error, got: %v", err)
	}
}

func TestReset_ClearsBaseline(t *testing.T) {
	f := &mockFetcher{result: map[string]string{"X": "old"}}
	p, _ := pin.New(f, pin.Config{MaxDrift: 0})

	p.Fetch(context.Background(), nil) // baseline = {X: old}

	f.result = map[string]string{"X": "new"}
	p.Reset()

	// After reset the new map becomes the baseline — no violation.
	_, err := p.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("expected no error after reset, got: %v", err)
	}
}
