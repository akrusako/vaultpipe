package dedupe_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/dedupe"
)

func TestFilter_EmptyInput(t *testing.T) {
	d := dedupe.New()
	out := d.Filter(map[string]string{})
	if len(out) != 0 {
		t.Fatalf("expected empty result, got %v", out)
	}
}

func TestFilter_AllNewKeys(t *testing.T) {
	d := dedupe.New()
	in := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out := d.Filter(in)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
}

func TestFilter_SuppressesUnchanged(t *testing.T) {
	d := dedupe.New()
	in := map[string]string{"FOO": "bar"}
	d.Filter(in) // prime the seen set
	out := d.Filter(in)
	if len(out) != 0 {
		t.Fatalf("expected empty result on second call, got %v", out)
	}
}

func TestFilter_PassesChangedValue(t *testing.T) {
	d := dedupe.New()
	d.Filter(map[string]string{"FOO": "old"})
	out := d.Filter(map[string]string{"FOO": "new"})
	if v, ok := out["FOO"]; !ok || v != "new" {
		t.Fatalf("expected FOO=new in result, got %v", out)
	}
}

func TestFilter_MixedChangedAndUnchanged(t *testing.T) {
	d := dedupe.New()
	d.Filter(map[string]string{"A": "1", "B": "2"})
	out := d.Filter(map[string]string{"A": "1", "B": "changed"})
	if _, ok := out["A"]; ok {
		t.Error("expected A to be suppressed")
	}
	if v, ok := out["B"]; !ok || v != "changed" {
		t.Errorf("expected B=changed, got %v", out)
	}
}

func TestReset_ClearsSeenSet(t *testing.T) {
	d := dedupe.New()
	d.Filter(map[string]string{"FOO": "bar"})
	d.Reset()
	out := d.Filter(map[string]string{"FOO": "bar"})
	if len(out) != 1 {
		t.Fatalf("expected 1 entry after reset, got %d", len(out))
	}
}

func TestSeen_ReturnsCopy(t *testing.T) {
	d := dedupe.New()
	d.Filter(map[string]string{"X": "y"})
	snap := d.Seen()
	snap["X"] = "mutated"
	if d.Seen()["X"] != "y" {
		t.Error("Seen should return a copy, not a reference")
	}
}
