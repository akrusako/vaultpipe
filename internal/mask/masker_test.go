package mask_test

import (
	"testing"

	"github.com/your-org/vaultpipe/internal/mask"
)

func TestNew_EmptySecrets(t *testing.T) {
	m := mask.New(map[string]string{})
	if m.Len() != 0 {
		t.Fatalf("expected 0 secrets, got %d", m.Len())
	}
}

func TestNew_IgnoresEmptyValues(t *testing.T) {
	m := mask.New(map[string]string{"KEY": ""})
	if m.Len() != 0 {
		t.Fatalf("expected empty values to be ignored")
	}
}

func TestRedact_ReplacesSecrets(t *testing.T) {
	m := mask.New(map[string]string{
		"DB_PASS": "supersecret",
		"API_KEY": "tok_abc123",
	})

	input := "connecting with supersecret and tok_abc123"
	want := "connecting with *** and ***"
	got := m.Redact(input)
	if got != want {
		t.Fatalf("Redact() = %q, want %q", got, want)
	}
}

func TestRedact_NoMatch(t *testing.T) {
	m := mask.New(map[string]string{"KEY": "hidden"})
	input := "nothing to see here"
	if got := m.Redact(input); got != input {
		t.Fatalf("Redact() modified string unexpectedly: %q", got)
	}
}

func TestAdd_RegistersNewSecret(t *testing.T) {
	m := mask.New(map[string]string{})
	m.Add("newtoken")
	if m.Len() != 1 {
		t.Fatalf("expected 1 secret after Add, got %d", m.Len())
	}
	if got := m.Redact("use newtoken here"); got != "use *** here" {
		t.Fatalf("Add secret not redacted: %q", got)
	}
}

func TestAdd_IgnoresEmpty(t *testing.T) {
	m := mask.New(map[string]string{})
	m.Add("")
	if m.Len() != 0 {
		t.Fatal("Add should ignore empty string")
	}
}
