package passthrough_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/passthrough"
)

func TestFilter_NoRules_PassesAll(t *testing.T) {
	p := passthrough.New()
	env := map[string]string{"HOME": "/root", "PATH": "/usr/bin", "SECRET": "s3cr3t"}
	out := p.Filter(env)
	if len(out) != len(env) {
		t.Fatalf("expected %d entries, got %d", len(env), len(out))
	}
}

func TestFilter_AllowPattern_KeepsMatches(t *testing.T) {
	p := passthrough.New(passthrough.WithAllow("path"))
	env := map[string]string{"PATH": "/usr/bin", "HOME": "/root"}
	out := p.Filter(env)
	if _, ok := out["PATH"]; !ok {
		t.Error("expected PATH to be present")
	}
	if _, ok := out["HOME"]; ok {
		t.Error("expected HOME to be absent")
	}
}

func TestFilter_DenyPattern_RemovesMatches(t *testing.T) {
	p := passthrough.New(passthrough.WithDeny("secret", "token", "vault"))
	env := map[string]string{
		"VAULT_TOKEN": "s.abc",
		"MY_SECRET":   "password",
		"HOME":        "/root",
	}
	out := p.Filter(env)
	if _, ok := out["VAULT_TOKEN"]; ok {
		t.Error("expected VAULT_TOKEN to be removed")
	}
	if _, ok := out["MY_SECRET"]; ok {
		t.Error("expected MY_SECRET to be removed")
	}
	if _, ok := out["HOME"]; !ok {
		t.Error("expected HOME to be present")
	}
}

func TestFilter_AllowAndDeny_DenyWins(t *testing.T) {
	p := passthrough.New(
		passthrough.WithAllow("vault"),
		passthrough.WithDeny("token"),
	)
	env := map[string]string{
		"VAULT_ADDR":  "https://vault:8200",
		"VAULT_TOKEN": "s.abc",
	}
	out := p.Filter(env)
	if _, ok := out["VAULT_ADDR"]; !ok {
		t.Error("expected VAULT_ADDR to be present")
	}
	if _, ok := out["VAULT_TOKEN"]; ok {
		t.Error("expected VAULT_TOKEN to be removed by deny rule")
	}
}

func TestFilter_CaseInsensitive(t *testing.T) {
	p := passthrough.New(passthrough.WithAllow("PATH"))
	env := map[string]string{"path": "/usr/bin", "HOME": "/root"}
	out := p.Filter(env)
	if _, ok := out["path"]; !ok {
		t.Error("expected lowercase 'path' to match allow pattern 'PATH'")
	}
}

func TestFilter_EmptyInput(t *testing.T) {
	p := passthrough.New(passthrough.WithAllow("home"))
	out := p.Filter(map[string]string{})
	if len(out) != 0 {
		t.Errorf("expected empty output, got %v", out)
	}
}

func TestFromOS_ReturnsNonEmpty(t *testing.T) {
	p := passthrough.New()
	out := p.FromOS()
	if len(out) == 0 {
		t.Error("expected at least one OS environment variable")
	}
}
