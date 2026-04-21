// Package scrub provides a Scrubber that removes sensitive or internal keys
// from a secret map before it is merged into a child process environment.
//
// By default the following keys (and any key sharing their prefix) are denied:
//
//   - VAULT_TOKEN  – the Vault authentication token
//   - VAULT_ADDR   – the Vault server address
//   - VAULTPIPE_   – all vaultpipe internal configuration variables
//
// Additional deny-list patterns can be registered at construction time via
// [WithPatterns]. Matching is case-insensitive and supports both exact-key and
// prefix-based rules.
package scrub
