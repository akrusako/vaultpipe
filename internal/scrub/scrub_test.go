package scrub_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/scrub"
)

func TestApply_PassesThroughCleanKeys(t *testing.T) {
	s := scrub.New()
	input := map[string]string{"DB_PASSWORD": "secret", "API_KEY": "abc"}
	out := s.Apply(input)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}

func TestApply_RemovesVaultToken(t *testing.T) {
	s := scrub.New()
	input := map[string]string{"VAULT_TOKEN": "hvs.xxx", "DB_PASS": "secret"}
	out := s.Apply(input)
	if _, ok := out["VAULT_TOKEN"]; ok {
		t.Fatal("VAULT_TOKEN should have been scrubbed")
	}
	if _, ok := out["DB_PASS"]; !ok {
		t.Fatal("DB_PASS should have been preserved")
	}
}

func TestApply_RemovesVaultAddr(t *testing.T) {
	s := scrub.New()
	input := map[string]string{"VAULT_ADDR": "http://127.0.0.1:8200"}
	out := s.Apply(input)
	if _, ok := out["VAULT_ADDR"]; ok {
		t.Fatal("VAULT_ADDR should have been scrubbed")
	}
}

func TestApply_RemovesVaultpipePrefix(t *testing.T) {
	s := scrub.New()
	input := map[string]string{
		"VAULTPIPE_SECRET_PATH": "/secret/foo",
		"VAULTPIPE_LOG_LEVEL":   "debug",
		"APP_ENV":               "production",
	}
	out := s.Apply(input)
	if len(out) != 1 {
		t.Fatalf("expected 1 key, got %d", len(out))
	}
	if _, ok := out["APP_ENV"]; !ok {
		t.Fatal("APP_ENV should have been preserved")
	}
}

func TestApply_CaseInsensitiveMatch(t *testing.T) {
	s := scrub.New()
	input := map[string]string{"vault_token": "hvs.lower"}
	out := s.Apply(input)
	if _, ok := out["vault_token"]; ok {
		t.Fatal("lowercase vault_token should have been scrubbed")
	}
}

func TestApply_CustomPattern(t *testing.T) {
	s := scrub.New(scrub.WithPatterns("INTERNAL_"))
	input := map[string]string{
		"INTERNAL_DEBUG": "true",
		"PUBLIC_KEY":     "rsa",
	}
	out := s.Apply(input)
	if _, ok := out["INTERNAL_DEBUG"]; ok {
		t.Fatal("INTERNAL_DEBUG should have been scrubbed by custom pattern")
	}
	if _, ok := out["PUBLIC_KEY"]; !ok {
		t.Fatal("PUBLIC_KEY should have been preserved")
	}
}

func TestApply_EmptyInput(t *testing.T) {
	s := scrub.New()
	out := s.Apply(map[string]string{})
	if len(out) != 0 {
		t.Fatalf("expected empty map, got %d keys", len(out))
	}
}
