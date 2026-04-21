package resolve_test

import (
	"os"
	"testing"

	"github.com/yourusername/vaultpipe/internal/resolve"
)

func TestExpand_NoVars(t *testing.T) {
	r := resolve.New(nil)
	got, err := r.Expand("secret/myapp/config")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "secret/myapp/config" {
		t.Errorf("got %q, want %q", got, "secret/myapp/config")
	}
}

func TestExpand_FromEnvMap(t *testing.T) {
	r := resolve.New(map[string]string{"ENV": "production"})
	got, err := r.Expand("secret/myapp/${ENV}/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "secret/myapp/production/db" {
		t.Errorf("got %q, want %q", got, "secret/myapp/production/db")
	}
}

func TestExpand_FallsBackToOS(t *testing.T) {
	os.Setenv("VP_REGION", "us-east-1")
	t.Cleanup(func() { os.Unsetenv("VP_REGION") })

	r := resolve.New(nil)
	got, err := r.Expand("secret/${VP_REGION}/certs")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "secret/us-east-1/certs" {
		t.Errorf("got %q, want %q", got, "secret/us-east-1/certs")
	}
}

func TestExpand_UndefinedVar(t *testing.T) {
	r := resolve.New(nil)
	_, err := r.Expand("secret/${UNDEFINED_XYZ_ABC}/key")
	if err == nil {
		t.Fatal("expected error for undefined variable, got nil")
	}
}

func TestExpand_EnvMapTakesPrecedenceOverOS(t *testing.T) {
	os.Setenv("APP", "from-os")
	t.Cleanup(func() { os.Unsetenv("APP") })

	r := resolve.New(map[string]string{"APP": "from-map"})
	got, err := r.Expand("secret/${APP}/token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "secret/from-map/token" {
		t.Errorf("got %q, want %q", got, "secret/from-map/token")
	}
}

func TestExpandAll_Success(t *testing.T) {
	r := resolve.New(map[string]string{"ENV": "staging"})
	paths := []string{"secret/${ENV}/db", "secret/${ENV}/cache"}
	got, err := r.ExpandAll(paths)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[0] != "secret/staging/db" || got[1] != "secret/staging/cache" {
		t.Errorf("unexpected results: %v", got)
	}
}

func TestExpandAll_StopsOnError(t *testing.T) {
	r := resolve.New(nil)
	paths := []string{"secret/static", "secret/${MISSING_VAR}/key"}
	_, err := r.ExpandAll(paths)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestHasVars(t *testing.T) {
	if resolve.HasVars("secret/plain/path") {
		t.Error("expected false for path without variables")
	}
	if !resolve.HasVars("secret/${ENV}/path") {
		t.Error("expected true for path with variable")
	}
}
