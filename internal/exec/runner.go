package exec

import (
	"fmt"
	"os"
	os_exec "os/exec"
	"strings"
)

// Runner holds the configuration for executing a subprocess with injected secrets.
type Runner struct {
	Env    map[string]string
	Stdout *os.File
	Stderr *os.File
}

// NewRunner creates a Runner with the provided secret environment map.
func NewRunner(env map[string]string) *Runner {
	return &Runner{
		Env:    env,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

// Run executes the given command with the merged environment (OS env + secrets).
// Secrets take precedence over existing OS environment variables.
func (r *Runner) Run(command string, args []string) error {
	if command == "" {
		return fmt.Errorf("exec: command must not be empty")
	}

	cmd := os_exec.Command(command, args...)
	cmd.Stdout = r.Stdout
	cmd.Stderr = r.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = r.buildEnv()

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("exec: process exited with error: %w", err)
	}
	return nil
}

// buildEnv merges OS environment variables with the injected secrets.
// Secret values override any existing OS variable with the same key.
func (r *Runner) buildEnv() []string {
	base := make(map[string]string, len(os.Environ()))
	for _, kv := range os.Environ() {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) == 2 {
			base[parts[0]] = parts[1]
		}
	}
	for k, v := range r.Env {
		base[k] = v
	}
	result := make([]string, 0, len(base))
	for k, v := range base {
		result = append(result, k+"="+v)
	}
	return result
}
