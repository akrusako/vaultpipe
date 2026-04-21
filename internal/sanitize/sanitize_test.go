package sanitize_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/sanitize"
)

func TestKey_Valid(t *testing.T) {
	s := sanitize.New()
	got, err := s.Key("MY_SECRET")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "MY_SECRET" {
		t.Errorf("got %q, want %q", got, "MY_SECRET")
	}
}

func TestKey_ReplacesHyphen(t *testing.T) {
	s := sanitize.New()
	got, err := s.Key("my-secret-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "my_secret_key" {
		t.Errorf("got %q, want %q", got, "my_secret_key")
	}
}

func TestKey_ReplacesDot(t *testing.T) {
	s := sanitize.New()
	got, err := s.Key("app.db.pass")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "app_db_pass" {
		t.Errorf("got %q, want %q", got, "app_db_pass")
	}
}

func TestKey_WithUppercase(t *testing.T) {
	s := sanitize.New(sanitize.WithUppercase())
	got, err := s.Key("my-secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "MY_SECRET" {
		t.Errorf("got %q, want %q", got, "MY_SECRET")
	}
}

func TestKey_InvalidStartsWithDigit(t *testing.T) {
	s := sanitize.New()
	_, err := s.Key("1bad")
	if err == nil {
		t.Fatal("expected error for key starting with digit")
	}
}

func TestKey_InvalidEmpty(t *testing.T) {
	s := sanitize.New()
	_, err := s.Key("")
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestMap_SkipsInvalidKeys(t *testing.T) {
	s := sanitize.New()
	in := map[string]string{
		"valid_key": "v1",
		"1invalid":  "v2",
	}
	out, errs := s.Map(in)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if _, ok := out["valid_key"]; !ok {
		t.Error("valid_key should be present in output")
	}
	if _, ok := out["1invalid"]; ok {
		t.Error("1invalid should not be present in output")
	}
}

func TestMap_NormalisesKeys(t *testing.T) {
	s := sanitize.New(sanitize.WithUppercase())
	in := map[string]string{"db-host": "localhost"}
	out, errs := s.Map(in)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %v", out)
	}
}
