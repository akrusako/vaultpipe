package label_test

import (
	"context"
	"errors"
	"testing"

	"github.com/yourusername/vaultpipe/internal/label"
)

// stubFetcher is a minimal Fetcher for tests.
type stubFetcher struct {
	data map[string]string
	err  error
}

func (s *stubFetcher) Fetch(_ context.Context, _ []string) (map[string]string, error) {
	if s.err != nil {
		return nil, s.err
	}
	out := make(map[string]string, len(s.data))
	for k, v := range s.data {
		out[k] = v
	}
	return out, nil
}

func TestNew_NilFetcherReturnsError(t *testing.T) {
	_, err := label.New(nil)
	if err == nil {
		t.Fatal("expected error for nil fetcher")
	}
}

func TestNew_NoOptions_PassThrough(t *testing.T) {
	stub := &stubFetcher{data: map[string]string{"KEY": "val"}}
	l, err := label.New(stub)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := l.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("fetch error: %v", err)
	}
	if got["KEY"] != "val" {
		t.Errorf("expected KEY=val, got %q", got["KEY"])
	}
}

func TestFetch_MergesLabels(t *testing.T) {
	stub := &stubFetcher{data: map[string]string{"SECRET": "s3cr3t"}}
	l, _ := label.New(stub,
		label.WithLabel("app", "myapp"),
		label.WithLabel("env", "production"),
	)
	got, err := l.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("fetch error: %v", err)
	}
	if got["app"] != "myapp" {
		t.Errorf("expected app=myapp, got %q", got["app"])
	}
	if got["env"] != "production" {
		t.Errorf("expected env=production, got %q", got["env"])
	}
	if got["SECRET"] != "s3cr3t" {
		t.Errorf("expected SECRET=s3cr3t, got %q", got["SECRET"])
	}
}

func TestFetch_SecretWinsOnKeyCollision(t *testing.T) {
	stub := &stubFetcher{data: map[string]string{"app": "from-vault"}}
	l, _ := label.New(stub, label.WithLabel("app", "from-label"))
	got, _ := l.Fetch(context.Background(), nil)
	if got["app"] != "from-vault" {
		t.Errorf("expected secret to win, got %q", got["app"])
	}
}

func TestFetch_PropagatesFetchError(t *testing.T) {
	stub := &stubFetcher{err: errors.New("vault unavailable")}
	l, _ := label.New(stub, label.WithLabel("env", "test"))
	_, err := l.Fetch(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error to propagate")
	}
}

func TestWithLabel_IgnoresEmptyKeyOrValue(t *testing.T) {
	stub := &stubFetcher{data: map[string]string{}}
	l, _ := label.New(stub,
		label.WithLabel("", "value"),
		label.WithLabel("key", ""),
		label.WithLabel("  ", "value"),
	)
	if len(l.Labels()) != 0 {
		t.Errorf("expected no labels, got %v", l.Labels())
	}
}

func TestLabels_ReturnsCopy(t *testing.T) {
	stub := &stubFetcher{data: map[string]string{}}
	l, _ := label.New(stub, label.WithLabel("k", "v"))
	copy1 := l.Labels()
	copy1["injected"] = "yes"
	if _, ok := l.Labels()["injected"]; ok {
		t.Error("Labels() should return an independent copy")
	}
}
