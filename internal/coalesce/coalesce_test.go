package coalesce_test

import (
	"testing"

	"github.com/your-org/vaultpipe/internal/coalesce"
)

func TestLookup_ReturnsFirstNonEmpty(t *testing.T) {
	c := coalesce.New(
		map[string]string{"KEY": ""},
		map[string]string{"KEY": "second"},
		map[string]string{"KEY": "third"},
	)
	v, err := c.Lookup("KEY")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "second" {
		t.Fatalf("expected %q, got %q", "second", v)
	}
}

func TestLookup_ErrNoValue_WhenAllEmpty(t *testing.T) {
	c := coalesce.New(
		map[string]string{"KEY": ""},
		map[string]string{},
	)
	_, err := c.Lookup("KEY")
	if err != coalesce.ErrNoValue {
		t.Fatalf("expected ErrNoValue, got %v", err)
	}
}

func TestLookup_ErrNoValue_WhenKeyAbsent(t *testing.T) {
	c := coalesce.New(map[string]string{"OTHER": "val"})
	_, err := c.Lookup("MISSING")
	if err != coalesce.ErrNoValue {
		t.Fatalf("expected ErrNoValue, got %v", err)
	}
}

func TestMerge_UnionOfAllKeys(t *testing.T) {
	c := coalesce.New(
		map[string]string{"A": "alpha", "B": ""},
		map[string]string{"B": "beta", "C": "gamma"},
	)
	out := c.Merge()
	if out["A"] != "alpha" {
		t.Errorf("A: expected %q, got %q", "alpha", out["A"])
	}
	if out["B"] != "beta" {
		t.Errorf("B: expected %q, got %q", "beta", out["B"])
	}
	if out["C"] != "gamma" {
		t.Errorf("C: expected %q, got %q", "gamma", out["C"])
	}
}

func TestMerge_EmptySourcesReturnsEmptyMap(t *testing.T) {
	c := coalesce.New()
	out := c.Merge()
	if len(out) != 0 {
		t.Fatalf("expected empty map, got %v", out)
	}
}

func TestAdd_AppendsLowerPrioritySource(t *testing.T) {
	c := coalesce.New(map[string]string{"X": "primary"})
	c.Add(map[string]string{"X": "fallback", "Y": "only-here"})

	if v, _ := c.Lookup("X"); v != "primary" {
		t.Errorf("X: expected %q, got %q", "primary", v)
	}
	if v, _ := c.Lookup("Y"); v != "only-here" {
		t.Errorf("Y: expected %q, got %q", "only-here", v)
	}
}
