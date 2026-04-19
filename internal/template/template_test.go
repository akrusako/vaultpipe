package template_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/template"
)

func TestRender_NoTemplate(t *testing.T) {
	r := template.New()
	out, err := r.Render("secret/static/path")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "secret/static/path" {
		t.Errorf("expected unchanged path, got %q", out)
	}
}

func TestRender_WithEnvVar(t *testing.T) {
	t.Setenv("VAULTPIPE_ENV", "production")
	r := template.New()
	out, err := r.Render("secret/{{.VAULTPIPE_ENV}}/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "secret/production/db" {
		t.Errorf("expected rendered path, got %q", out)
	}
}

func TestRender_MissingKey(t *testing.T) {
	r := template.New()
	_, err := r.Render("secret/{{.DOES_NOT_EXIST}}/db")
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestRender_InvalidTemplate(t *testing.T) {
	r := template.New()
	_, err := r.Render("secret/{{.UNCLOSED")
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
}

func TestRender_EmptyString(t *testing.T) {
	r := template.New()
	out, err := r.Render("")
	if err != nil {
		t.Fatalf("unexpected error for empty string: %v", err)
	}
	if out != "" {
		t.Errorf("expected empty string, got %q", out)
	}
}

func TestRenderAll_Mixed(t *testing.T) {
	t.Setenv("VAULTPIPE_REGION", "us-east-1")
	r := template.New()
	paths := []string{
		"secret/static",
		"secret/{{.VAULTPIPE_REGION}}/app",
	}
	out, err := r.RenderAll(paths)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 paths, got %d", len(out))
	}
	if out[1] != "secret/us-east-1/app" {
		t.Errorf("expected rendered path, got %q", out[1])
	}
}

func TestRenderAll_ErrorPropagates(t *testing.T) {
	r := template.New()
	paths := []string{"secret/static", "secret/{{.MISSING_KEY}}/app"}
	_, err := r.RenderAll(paths)
	if err == nil {
		t.Fatal("expected error to propagate, got nil")
	}
}
