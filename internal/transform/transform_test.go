package transform_test

import (
	"errors"
	"testing"

	"github.com/yourusername/vaultpipe/internal/transform"
)

func TestApply_NoFuncs(t *testing.T) {
	tr := transform.New()
	in := map[string]string{"foo": "bar"}
	out, err := tr.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["foo"] != "bar" {
		t.Errorf("expected bar, got %q", out["foo"])
	}
}

func TestApply_PrefixKeys(t *testing.T) {
	tr := transform.New(transform.PrefixKeys("APP_"))
	out, err := tr.Apply(map[string]string{"db_pass": "secret"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["APP_db_pass"]; !ok {
		t.Errorf("expected key APP_db_pass, got %v", out)
	}
}

func TestApply_UppercaseKeys(t *testing.T) {
	tr := transform.New(transform.UppercaseKeys())
	out, err := tr.Apply(map[string]string{"api_key": "xyz"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["API_KEY"] != "xyz" {
		t.Errorf("expected API_KEY=xyz, got %v", out)
	}
}

func TestApply_TrimValueSpace(t *testing.T) {
	tr := transform.New(transform.TrimValueSpace())
	out, err := tr.Apply(map[string]string{"token": "  abc123  "})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["token"] != "abc123" {
		t.Errorf("expected trimmed value, got %q", out["token"])
	}
}

func TestApply_ChainedFuncs(t *testing.T) {
	tr := transform.New(transform.PrefixKeys("X_"), transform.UppercaseKeys())
	out, err := tr.Apply(map[string]string{"secret": "val"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["X_SECRET"] != "val" {
		t.Errorf("expected X_SECRET=val, got %v", out)
	}
}

func TestApply_FuncReturnsError(t *testing.T) {
	failFn := func(k, v string) (string, string, error) {
		return "", "", errors.New("boom")
	}
	tr := transform.New(failFn)
	_, err := tr.Apply(map[string]string{"k": "v"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestApply_EmptySecrets(t *testing.T) {
	tr := transform.New(transform.UppercaseKeys())
	out, err := tr.Apply(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}
