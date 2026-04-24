// Package passthrough controls which host environment variables are forwarded
// to child processes launched by vaultpipe.
//
// By default all host variables are forwarded. Callers may restrict the set
// with allow patterns (only matching keys are forwarded) and further narrow it
// with deny patterns (matching keys are always removed). Both pattern lists
// are evaluated case-insensitively using substring matching.
//
// Typical usage:
//
//	pt := passthrough.New(
//		passthrough.WithDeny("vault_token", "vault_addr", "vaultpipe"),
//	)
//	env := pt.FromOS()
//
// The Middleware type integrates passthrough into the secret-fetch pipeline so
// that host variables and Vault secrets are merged in a single step, with Vault
// values taking precedence.
package passthrough
