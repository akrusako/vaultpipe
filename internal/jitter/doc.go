// Package jitter provides a Spreader that introduces a randomised offset to
// a base interval. It is used by the secret-watch and lease-renewal loops to
// prevent thundering-herd behaviour when many vaultpipe processes share the
// same Vault cluster.
//
// Usage:
//
//	s, err := jitter.New(jitter.Config{Base: 30 * time.Second, Factor: 0.2})
//	if err != nil { ... }
//	// blocks for 30s–36s then returns
//	if err := s.Wait(ctx); err != nil { ... }
package jitter
