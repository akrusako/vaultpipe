package ceiling_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/yourusername/vaultpipe/internal/ceiling"
)

func makeFetcher(secrets map[string]string, err error) ceiling.FetchFunc {
	return func(_ context.Context, _ []string) (map[string]string, error) {
		return secrets, err
	}
}

func TestNew_DefaultMax(t *testing.T) {
	l := ceiling.New(makeFetcher(nil, nil), 0)
	if l.Max() != ceiling.DefaultMax {
		t.Fatalf("expected default max %d, got %d", ceiling.DefaultMax, l.Max())
	}
}

func TestNew_CustomMax(t *testing.T) {
	l := ceiling.New(makeFetcher(nil, nil), 10)
	if l.Max() != 10 {
		t.Fatalf("expected max 10, got %d", l.Max())
	}
}

func TestFetch_UnderLimit(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2"}
	l := ceiling.New(makeFetcher(secrets, nil), 5)
	got, err := l.Fetch(context.Background(), []string{"secret/data/app"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(got))
	}
}

func TestFetch_AtLimit(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2", "C": "3"}
	l := ceiling.New(makeFetcher(secrets, nil), 3)
	_, err := l.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("expected no error at exact limit, got: %v", err)
	}
}

func TestFetch_ExceedsLimit(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2", "C": "3", "D": "4"}
	l := ceiling.New(makeFetcher(secrets, nil), 3)
	_, err := l.Fetch(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error when limit exceeded, got nil")
	}
	if !strings.Contains(err.Error(), "ceiling") {
		t.Fatalf("expected ceiling error, got: %v", err)
	}
}

func TestFetch_PropagatesUpstreamError(t *testing.T) {
	upstreamErr := errors.New("vault unavailable")
	l := ceiling.New(makeFetcher(nil, upstreamErr), 10)
	_, err := l.Fetch(context.Background(), nil)
	if !errors.Is(err, upstreamErr) {
		t.Fatalf("expected upstream error, got: %v", err)
	}
}
