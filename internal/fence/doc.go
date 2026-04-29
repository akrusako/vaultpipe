// Package fence implements a write-once injection-window guard for vaultpipe.
//
// A Fence ensures that the secret injection window — the period between
// reading secrets from Vault and handing them to the child process — is
// entered and exited exactly once per execution cycle.
//
// Usage:
//
//	f := fence.New(fence.WithTTL(5 * time.Second))
//
//	token := uuid.New().String()
//	if err := f.Open(token); err != nil {
//		// already in an injection window
//	}
//
//	// ... inject secrets into child process environment ...
//
//	if err := f.Close(token); err != nil {
//		// window expired or token mismatch
//	}
package fence
