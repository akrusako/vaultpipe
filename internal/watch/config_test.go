package watch_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/watch"
)

func TestDefaultConfig(t *testing.T) {
	cfg := watch.DefaultConfig()
	if cfg.Interval != 30*time.Second {
		t.Errorf("expected 30s interval, got %v", cfg.Interval)
	}
}

func TestValidate_MissingPaths(t *testing.T) {
	cfg := watch.Config{Interval: 10 * time.Second}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing paths")
	}
}

func TestValidate_ZeroInterval(t *testing.T) {
	cfg := watch.Config{Paths: []string{"secret/app"}, Interval: 0}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for zero interval")
	}
}

func TestValidate_Valid(t *testing.T) {
	cfg := watch.Config{Paths: []string{"secret/app"}, Interval: 5 * time.Second}
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNewFromConfig_InvalidConfig(t *testing.T) {
	_, err := watch.NewFromConfig(watch.Config{}, nil, nil)
	if err == nil {
		t.Error("expected error for invalid config")
	}
}

func TestNewFromConfig_Valid(t *testing.T) {
	cfg := watch.Config{
		Paths:    []string{"secret/app"},
		Interval: 15 * time.Second,
	}
	fetcher := &mockFetcher{results: []map[string]string{{"K": "v"}}}
	w, err := watch.NewFromConfig(cfg, fetcher, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w == nil {
		t.Error("expected non-nil Watcher")
	}
}
