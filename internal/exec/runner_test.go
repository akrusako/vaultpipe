package exec

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestNewRunner(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	r := NewRunner(env)
	if r.Env["FOO"] != "bar" {
		t.Fatalf("expected FOO=bar, got %s", r.Env["FOO"])
	}
}

func TestRun_EmptyCommand(t *testing.T) {
	r := NewRunner(nil)
	err := r.Run("", nil)
	if err == nil || !strings.Contains(err.Error(), "command must not be empty") {
		t.Fatalf("expected empty command error, got %v", err)
	}
}

func TestRun_Success(t *testing.T) {
	r := NewRunner(map[string]string{"VAULTPIPE_TEST": "injected"})
	// Capture stdout by redirecting runner output
	pr, pw, _ := os.Pipe()
	r.Stdout = pw
	r.Stderr = pw

	err := r.Run("sh", []string{"-c", "echo $VAULTPIPE_TEST"})
	pw.Close()

	var buf bytes.Buffer
	buf.ReadFrom(pr)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "injected") {
		t.Fatalf("expected 'injected' in output, got: %s", buf.String())
	}
}

func TestRun_NonZeroExit(t *testing.T) {
	r := NewRunner(nil)
	err := r.Run("sh", []string{"-c", "exit 1"})
	if err == nil {
		t.Fatal("expected error for non-zero exit")
	}
}

func TestBuildEnv_SecretOverridesOS(t *testing.T) {
	os.Setenv("OVERRIDE_ME", "original")
	defer os.Unsetenv("OVERRIDE_ME")

	r := NewRunner(map[string]string{"OVERRIDE_ME": "secret_value"})
	env := r.buildEnv()

	for _, kv := range env {
		if strings.HasPrefix(kv, "OVERRIDE_ME=") {
			if kv != "OVERRIDE_ME=secret_value" {
				t.Fatalf("expected secret override, got %s", kv)
			}
			return
		}
	}
	t.Fatal("OVERRIDE_ME not found in built env")
}
