package env

import (
	"sort"
	"testing"
)

func TestNew_CopiesBase(t *testing.T) {
	base := map[string]string{"HOME": "/root", "PATH": "/usr/bin"}
	w := New(base)
	base["HOME"] = "mutated"
	if v, _ := w.Lookup("HOME"); v != "/root" {
		t.Fatalf("expected /root, got %s", v)
	}
}

func TestApply_OverlaysSecrets(t *testing.T) {
	w := New(map[string]string{"EXISTING": "old"})
	w.Apply(map[string]string{"DB_PASS": "secret", "EXISTING": "new"})

	if v, ok := w.Lookup("DB_PASS"); !ok || v != "secret" {
		t.Fatalf("expected DB_PASS=secret, got %s", v)
	}
	if v, _ := w.Lookup("EXISTING"); v != "new" {
		t.Fatalf("overlay should shadow base: got %s", v)
	}
}

func TestBuild_MergesAll(t *testing.T) {
	w := New(map[string]string{"A": "1", "B": "2"})
	w.Apply(map[string]string{"B": "override", "C": "3"})

	env := w.Build()
	m := ParseEnvSlice(env)

	if m["A"] != "1" {
		t.Errorf("A: want 1, got %s", m["A"])
	}
	if m["B"] != "override" {
		t.Errorf("B: want override, got %s", m["B"])
	}
	if m["C"] != "3" {
		t.Errorf("C: want 3, got %s", m["C"])
	}
}

func TestBuild_EmptyBase(t *testing.T) {
	w := New(nil)
	w.Apply(map[string]string{"TOKEN": "abc"})
	env := w.Build()
	if len(env) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(env))
	}
}

func TestLookup_MissingKey(t *testing.T) {
	w := New(map[string]string{})
	if _, ok := w.Lookup("NOPE"); ok {
		t.Fatal("expected missing key to return false")
	}
}

func TestParseEnvSlice_Standard(t *testing.T) {
	slice := []string{"FOO=bar", "BAZ=qux=extra", "NOVALUE"}
	m := ParseEnvSlice(slice)

	cases := map[string]string{"FOO": "bar", "BAZ": "qux=extra", "NOVALUE": ""}
	for k, want := range cases {
		if got := m[k]; got != want {
			t.Errorf("%s: want %q, got %q", k, want, got)
		}
	}
}

func TestBuild_Deterministic(t *testing.T) {
	w := New(map[string]string{"Z": "26", "A": "1"})
	env := w.Build()
	sorted := make([]string, len(env))
	copy(sorted, env)
	sort.Strings(sorted)
	if len(sorted) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(sorted))
	}
}
