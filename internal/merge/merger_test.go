package merge_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/merge"
)

func TestMerge_EmptyInput(t *testing.T) {
	m := merge.New(merge.StrategyLast)
	out, err := m.Merge()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("expected empty map, got %v", out)
	}
}

func TestMerge_SingleMap(t *testing.T) {
	m := merge.New(merge.StrategyFirst)
	out, err := m.Merge(map[string]string{"A": "1", "B": "2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "1" || out["B"] != "2" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestMerge_StrategyFirst_KeepsEarlierValue(t *testing.T) {
	m := merge.New(merge.StrategyFirst)
	a := map[string]string{"KEY": "first"}
	b := map[string]string{"KEY": "second"}
	out, err := m.Merge(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "first" {
		t.Fatalf("expected 'first', got %q", out["KEY"])
	}
}

func TestMerge_StrategyLast_KeepsLaterValue(t *testing.T) {
	m := merge.New(merge.StrategyLast)
	a := map[string]string{"KEY": "first"}
	b := map[string]string{"KEY": "second"}
	out, err := m.Merge(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "second" {
		t.Fatalf("expected 'second', got %q", out["KEY"])
	}
}

func TestMerge_StrategyError_ReturnsDuplicateError(t *testing.T) {
	m := merge.New(merge.StrategyError)
	a := map[string]string{"KEY": "first"}
	b := map[string]string{"KEY": "second"}
	_, err := m.Merge(a, b)
	if err == nil {
		t.Fatal("expected error for duplicate key, got nil")
	}
}

func TestMerge_DoesNotMutateInputs(t *testing.T) {
	m := merge.New(merge.StrategyLast)
	a := map[string]string{"A": "1"}
	b := map[string]string{"B": "2"}
	out, _ := m.Merge(a, b)
	out["A"] = "mutated"
	if a["A"] != "1" {
		t.Fatal("input map was mutated")
	}
}

func TestMustMerge_PanicsOnDuplicate(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic, got none")
		}
	}()
	m := merge.New(merge.StrategyError)
	m.MustMerge(
		map[string]string{"X": "1"},
		map[string]string{"X": "2"},
	)
}
