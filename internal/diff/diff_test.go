package diff

import (
	"testing"
)

func TestCompare_NoChanges(t *testing.T) {
	d := New()
	prev := map[string]string{"FOO": "bar", "BAZ": "qux"}
	next := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if changes := d.Compare(prev, next); len(changes) != 0 {
		t.Fatalf("expected no changes, got %d", len(changes))
	}
}

func TestCompare_Added(t *testing.T) {
	d := New()
	prev := map[string]string{}
	next := map[string]string{"NEW_KEY": "value"}
	changes := d.Compare(prev, next)
	if len(changes) != 1 || changes[0].Action != "added" || changes[0].Key != "NEW_KEY" {
		t.Fatalf("expected one added change, got %+v", changes)
	}
}

func TestCompare_Removed(t *testing.T) {
	d := New()
	prev := map[string]string{"OLD_KEY": "value"}
	next := map[string]string{}
	changes := d.Compare(prev, next)
	if len(changes) != 1 || changes[0].Action != "removed" || changes[0].Key != "OLD_KEY" {
		t.Fatalf("expected one removed change, got %+v", changes)
	}
}

func TestCompare_Changed(t *testing.T) {
	d := New()
	prev := map[string]string{"TOKEN": "old"}
	next := map[string]string{"TOKEN": "new"}
	changes := d.Compare(prev, next)
	if len(changes) != 1 || changes[0].Action != "changed" {
		t.Fatalf("expected one changed entry, got %+v", changes)
	}
}

func TestHasChanges_True(t *testing.T) {
	d := New()
	if !d.HasChanges(map[string]string{"A": "1"}, map[string]string{"A": "2"}) {
		t.Fatal("expected HasChanges to return true")
	}
}

func TestHasChanges_False(t *testing.T) {
	d := New()
	m := map[string]string{"A": "1"}
	if d.HasChanges(m, m) {
		t.Fatal("expected HasChanges to return false")
	}
}
