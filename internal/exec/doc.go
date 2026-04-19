// Package exec provides utilities for spawning subprocesses with secrets
// injected directly into the process environment.
//
// Secrets sourced from Vault are merged with the current OS environment and
// passed to the child process. Secret values take precedence over any
// pre-existing environment variable with the same key. No secret data is
// written to disk at any point during execution.
//
// Basic usage:
//
//	env, _ := vaultClient.ReadSecrets("secret/myapp")
//	runner := exec.NewRunner(env)
//	runner.Run("myapp", []string{"--serve"})
package exec
