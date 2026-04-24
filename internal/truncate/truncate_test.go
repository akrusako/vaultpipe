package truncate_test

import (
	"strings"
	"testing"

	"github.com/yourusername/vaultpipe/internal/truncate"
)

func TestNew_DefaultMaxBytes(t *testing.T) {
	tr := truncate.New()
	if tr.MaxBytes() != truncate.DefaultMaxBytes {
		t.Fatalf("expected default %d, got %d", truncate.DefaultMaxBytes, tr.MaxBytes())
	}
}

func TestWithMaxBytes_SetsValue(t *testing.T) {
	tr := truncate.New(truncate.WithMaxBytes(64))
	if tr.MaxBytes() != 64 {
		t.Fatalf("expected 64, got %d", tr.MaxBytes())
	}
}

func TestWithMaxBytes_PanicsOnZero(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero maxBytes")
		}
	}()
	truncate.WithMaxBytes(0)
}

func TestWithMaxBytes_PanicsOnNegative(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for negative maxBytes")
		}
	}()
	truncate.WithMaxBytes(-1)
}

func TestApply_ShortValuesUnchanged(t *testing.T) {
	tr := truncate.New(truncate.WithMaxBytes(10))
	in := map[string]string{"KEY": "hello"}
	out := tr.Apply(in)
	if out["KEY"] != "hello" {
		t.Fatalf("expected 'hello', got %q", out["KEY"])
	}
}

func TestApply_TruncatesLongValue(t *testing.T) {
	const max = 8
	tr := truncate.New(truncate.WithMaxBytes(max))
	long := strings.Repeat("x", 100)
	in := map[string]string{"TOKEN": long}
	out := tr.Apply(in)
	if len(out["TOKEN"]) != max {
		t.Fatalf("expected length %d, got %d", max, len(out["TOKEN"]))
	}
}

func TestApply_DoesNotModifyOriginal(t *testing.T) {
	tr := truncate.New(truncate.WithMaxBytes(4))
	orig := "abcdefgh"
	in := map[string]string{"K": orig}
	_ = tr.Apply(in)
	if in["K"] != orig {
		t.Fatalf("original map was mutated")
	}
}

func TestApply_EmptyMap(t *testing.T) {
	tr := truncate.New()
	out := tr.Apply(map[string]string{})
	if len(out) != 0 {
		t.Fatalf("expected empty map, got %v", out)
	}
}

func TestApply_ExactLengthUnchanged(t *testing.T) {
	const max = 5
	tr := truncate.New(truncate.WithMaxBytes(max))
	val := "abcde" // exactly max bytes
	out := tr.Apply(map[string]string{"K": val})
	if out["K"] != val {
		t.Fatalf("expected %q unchanged, got %q", val, out["K"])
	}
}
