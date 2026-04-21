// Package envelope provides AES-GCM envelope encryption for secret values
// managed by vaultpipe.
//
// Secrets fetched from Vault can be sealed before being stored in memory
// or passed across process boundaries, and opened only when injected into
// the child process environment. This limits the window during which
// plaintext secret values are resident in memory.
//
// Usage:
//
//	e, err := envelope.New(key) // key must be 16, 24, or 32 bytes
//	sealed, err := e.Seal("my-secret")
//	plain, err  := e.Open(sealed)
//
// For bulk operations over a map of secrets use SealMap and OpenMap.
package envelope
