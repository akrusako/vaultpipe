package redact_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/redact"
)

func TestNew_FiltersEmptySecrets(t *testing.T) {
	r := redact.New("***", []string{"", "  ", "valid"})
	if r.Len() != 1 {
		t.Fatalf("expected 1 secret, got %d", r.Len())
	}
}

func TestScrub_ReplacesSecret(t *testing.T) {
	r := redact.New("[REDACTED]", []string{"s3cr3t"})
	out := r.Scrub("the password is s3cr3t ok")
	if out != "the password is [REDACTED] ok" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestScrub_MultipleSecrets(t *testing.T) {
	r := redact.New("***", []string{"alpha", "beta"})
	out := r.Scrub("alpha and beta are secrets")
	if out != "*** and *** are secrets" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestScrub_NoMatch(t *testing.T) {
	r := redact.New("***", []string{"missing"})
	input := "nothing to replace here"
	if got := r.Scrub(input); got != input {
		t.Fatalf("expected unchanged string, got %q", got)
	}
}

func TestAdd_RegistersSecret(t *testing.T) {
	r := redact.New("***", nil)
	r.Add("newtoken")
	if r.Len() != 1 {
		t.Fatalf("expected 1 secret after Add, got %d", r.Len())
	}
	out := r.Scrub("token=newtoken")
	if out != "token=***" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestAdd_IgnoresEmpty(t *testing.T) {
	r := redact.New("***", nil)
	r.Add("")
	r.Add("   ")
	if r.Len() != 0 {
		t.Fatalf("expected 0 secrets, got %d", r.Len())
	}
}
