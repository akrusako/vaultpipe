package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/checkpoint"
)

func tempFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "checkpoint.json")
}

func TestNew_NonExistentFile(t *testing.T) {
	cp, err := checkpoint.New(tempFile(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cp.Has("secret/data/foo") {
		t.Fatal("expected empty checkpoint")
	}
}

func TestMark_PersistsRecord(t *testing.T) {
	path := tempFile(t)
	cp, _ := checkpoint.New(path)

	if err := cp.Mark("secret/data/db"); err != nil {
		t.Fatalf("Mark failed: %v", err)
	}

	if !cp.Has("secret/data/db") {
		t.Fatal("expected Has to return true after Mark")
	}
}

func TestMark_WritesToDisk(t *testing.T) {
	path := tempFile(t)
	cp, _ := checkpoint.New(path)
	_ = cp.Mark("secret/data/api")

	// Re-load from disk
	cp2, err := checkpoint.New(path)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	if !cp2.Has("secret/data/api") {
		t.Fatal("expected reloaded checkpoint to contain path")
	}
}

func TestGet_ReturnsRecord(t *testing.T) {
	path := tempFile(t)
	cp, _ := checkpoint.New(path)
	before := time.Now().UTC()
	_ = cp.Mark("secret/data/tls")

	rec, ok := cp.Get("secret/data/tls")
	if !ok {
		t.Fatal("expected record to be found")
	}
	if rec.Path != "secret/data/tls" {
		t.Errorf("got path %q, want %q", rec.Path, "secret/data/tls")
	}
	if rec.FetchedAt.Before(before) {
		t.Error("FetchedAt should be >= time before Mark")
	}
}

func TestGet_Missing(t *testing.T) {
	cp, _ := checkpoint.New(tempFile(t))
	_, ok := cp.Get("secret/data/missing")
	if ok {
		t.Fatal("expected missing record")
	}
}

func TestReset_ClearsRecords(t *testing.T) {
	path := tempFile(t)
	cp, _ := checkpoint.New(path)
	_ = cp.Mark("secret/data/foo")
	_ = cp.Mark("secret/data/bar")

	if err := cp.Reset(); err != nil {
		t.Fatalf("Reset failed: %v", err)
	}
	if cp.Has("secret/data/foo") || cp.Has("secret/data/bar") {
		t.Fatal("expected empty checkpoint after Reset")
	}

	// File should still exist and be loadable
	cp2, err := checkpoint.New(path)
	if err != nil {
		t.Fatalf("reload after Reset failed: %v", err)
	}
	if cp2.Has("secret/data/foo") {
		t.Fatal("expected reloaded checkpoint to be empty after Reset")
	}
}

func TestNew_InvalidJSON(t *testing.T) {
	path := tempFile(t)
	_ = os.WriteFile(path, []byte("not json"), 0o600)
	_, err := checkpoint.New(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
