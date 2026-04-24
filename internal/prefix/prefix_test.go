package prefix_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/prefix"
)

func TestApply_EmptyPrefix(t *testing.T) {
	pr := prefix.New()
	input := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out := pr.Apply(input)
	if len(out) != 2 || out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Fatalf("expected original map unchanged, got %v", out)
	}
}

func TestApply_AddsPrefix(t *testing.T) {
	pr := prefix.New(prefix.WithPrefix("APP_"))
	input := map[string]string{"SECRET": "s3cr3t", "TOKEN": "tok"}
	out := pr.Apply(input)
	if out["APP_SECRET"] != "s3cr3t" {
		t.Errorf("expected APP_SECRET=s3cr3t, got %q", out["APP_SECRET"])
	}
	if out["APP_TOKEN"] != "tok" {
		t.Errorf("expected APP_TOKEN=tok, got %q", out["APP_TOKEN"])
	}
	if len(out) != 2 {
		t.Errorf("unexpected extra keys: %v", out)
	}
}

func TestApply_TrimsPrefixWhitespace(t *testing.T) {
	pr := prefix.New(prefix.WithPrefix("  VAULT_  "))
	out := pr.Apply(map[string]string{"KEY": "val"})
	if _, ok := out["VAULT_KEY"]; !ok {
		t.Fatalf("expected VAULT_KEY, got %v", out)
	}
}

func TestApply_EmptyInput(t *testing.T) {
	pr := prefix.New(prefix.WithPrefix("X_"))
	out := pr.Apply(map[string]string{})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestStrip_RemovesPrefixedKeys(t *testing.T) {
	pr := prefix.New(prefix.WithPrefix("APP_"))
	input := map[string]string{
		"APP_SECRET": "s3cr3t",
		"APP_TOKEN":  "tok",
		"OTHER":      "ignored",
	}
	out := pr.Strip(input)
	if out["SECRET"] != "s3cr3t" {
		t.Errorf("expected SECRET=s3cr3t, got %q", out["SECRET"])
	}
	if out["TOKEN"] != "tok" {
		t.Errorf("expected TOKEN=tok, got %q", out["TOKEN"])
	}
	if _, ok := out["OTHER"]; ok {
		t.Error("expected OTHER to be dropped")
	}
}

func TestStrip_EmptyPrefix_CopiesAll(t *testing.T) {
	pr := prefix.New()
	input := map[string]string{"A": "1", "B": "2"}
	out := pr.Strip(input)
	if len(out) != 2 || out["A"] != "1" || out["B"] != "2" {
		t.Fatalf("expected full copy, got %v", out)
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	pr := prefix.New(prefix.WithPrefix("P_"))
	input := map[string]string{"KEY": "val"}
	pr.Apply(input)
	if _, ok := input["P_KEY"]; ok {
		t.Error("Apply must not mutate the input map")
	}
}
