package env_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/yourusername/vaultpipe/internal/env"
)

// TestIntegration_EnvReachesChildProcess verifies that secrets injected via
// Writer.Build() are visible inside a real child process.
func TestIntegration_EnvReachesChildProcess(t *testing.T) {
	if os.Getenv("CI") == "" {
		t.Skip("integration test; set CI=1 to run")
	}

	base := env.ParseEnvSlice(os.Environ())
	w := env.New(base)
	w.Apply(map[string]string{
		"VAULTPIPE_TEST_SECRET": "hunter2",
	})

	cmd := exec.Command("sh", "-c", "echo $VAULTPIPE_TEST_SECRET")
	cmd.Env = w.Build()
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("child process error: %v", err)
	}
	got := strings.TrimSpace(string(out))
	if got != "hunter2" {
		t.Errorf("expected hunter2, got %q", got)
	}
}

// TestIntegration_BaseEnvPreserved ensures the original OS environment
// remains accessible in the child when no overlay is provided.
func TestIntegration_BaseEnvPreserved(t *testing.T) {
	if os.Getenv("CI") == "" {
		t.Skip("integration test; set CI=1 to run")
	}

	t.Setenv("VAULTPIPE_PRESERVED", "yes")

	base := env.ParseEnvSlice(os.Environ())
	w := env.New(base)

	cmd := exec.Command("sh", "-c", "echo $VAULTPIPE_PRESERVED")
	cmd.Env = w.Build()
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("child process error: %v", err)
	}
	if got := strings.TrimSpace(string(out)); got != "yes" {
		t.Errorf("expected yes, got %q", got)
	}
}
