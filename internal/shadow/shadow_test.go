package shadow_test

import (
	"errors"
	"testing"

	"github.com/yourusername/vaultpipe/internal/shadow"
)

func TestSet_AndGet_RoundTrip(t *testing.T) {
	s := shadow.New()
	if err := s.Set("KEY", "value"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := s.Get("KEY")
	if !ok {
		t.Fatal("expected key to be present")
	}
	if v != "value" {
		t.Fatalf("got %q, want %q", v, "value")
	}
}

func TestSet_RejectsOverwrite(t *testing.T) {
	s := shadow.New()
	_ = s.Set("KEY", "first")
	err := s.Set("KEY", "second")
	if err == nil {
		t.Fatal("expected ErrKeyExists, got nil")
	}
	if !errors.Is(err, shadow.ErrKeyExists) {
		t.Fatalf("expected ErrKeyExists, got %v", err)
	}
	// Original value must be preserved.
	v, _ := s.Get("KEY")
	if v != "first" {
		t.Fatalf("value was mutated: got %q", v)
	}
}

func TestGet_MissingKey(t *testing.T) {
	s := shadow.New()
	_, ok := s.Get("MISSING")
	if ok {
		t.Fatal("expected missing key to return false")
	}
}

func TestSnapshot_IsCopy(t *testing.T) {
	s := shadow.New()
	_ = s.Set("A", "1")
	snap := s.Snapshot()
	snap["A"] = "mutated"
	v, _ := s.Get("A")
	if v != "1" {
		t.Fatal("snapshot mutation affected store")
	}
}

func TestLen_ReflectsEntries(t *testing.T) {
	s := shadow.New()
	if s.Len() != 0 {
		t.Fatal("expected empty store")
	}
	_ = s.Set("X", "a")
	_ = s.Set("Y", "b")
	if s.Len() != 2 {
		t.Fatalf("expected 2, got %d", s.Len())
	}
}

func TestLoadMap_BulkInsert(t *testing.T) {
	s := shadow.New()
	m := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if err := s.LoadMap(m); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", s.Len())
	}
}

func TestLoadMap_StopsOnCollision(t *testing.T) {
	s := shadow.New()
	_ = s.Set("FOO", "original")
	err := s.LoadMap(map[string]string{"FOO": "new"})
	if !errors.Is(err, shadow.ErrKeyExists) {
		t.Fatalf("expected ErrKeyExists, got %v", err)
	}
}
