package snapshot

import (
	"testing"
)

func TestLatest_NilBeforeSave(t *testing.T) {
	s := New()
	if s.Latest() != nil {
		t.Fatal("expected nil before any save")
	}
}

func TestSave_StoresSecrets(t *testing.T) {
	s := New()
	secrets := map[string]string{"DB_PASS": "hunter2"}
	s.Save(secrets)
	snap := s.Latest()
	if snap == nil {
		t.Fatal("expected snapshot after save")
	}
	if snap.Secrets["DB_PASS"] != "hunter2" {
		t.Errorf("unexpected value: %s", snap.Secrets["DB_PASS"])
	}
}

func TestSave_CopiesMap(t *testing.T) {
	s := New()
	secrets := map[string]string{"KEY": "val"}
	s.Save(secrets)
	secrets["KEY"] = "mutated"
	if s.Latest().Secrets["KEY"] != "val" {
		t.Error("snapshot should not reflect mutation of original map")
	}
}

func TestChanged_TrueWhenNoSnapshot(t *testing.T) {
	s := New()
	if !s.Changed(map[string]string{"A": "1"}) {
		t.Error("expected changed=true with no prior snapshot")
	}
}

func TestChanged_FalseWhenIdentical(t *testing.T) {
	s := New()
	secrets := map[string]string{"A": "1", "B": "2"}
	s.Save(secrets)
	if s.Changed(map[string]string{"A": "1", "B": "2"}) {
		t.Error("expected changed=false for identical secrets")
	}
}

func TestChanged_TrueOnValueChange(t *testing.T) {
	s := New()
	s.Save(map[string]string{"A": "old"})
	if !s.Changed(map[string]string{"A": "new"}) {
		t.Error("expected changed=true when value differs")
	}
}

func TestChanged_TrueOnAddedKey(t *testing.T) {
	s := New()
	s.Save(map[string]string{"A": "1"})
	if !s.Changed(map[string]string{"A": "1", "B": "2"}) {
		t.Error("expected changed=true when key is added")
	}
}
