// Package renew provides automatic Vault token renewal for long-running
// child processes managed by vaultpipe.
//
// A Renewer is created with a Vault API client and a renewal interval.
// Calling Start blocks and renews the token on each tick until the
// provided context is cancelled — typically when the child process exits.
//
// Usage:
//
//	r := renew.New(vaultClient, 5*time.Minute, logger)
//	go r.Start(ctx)
package renew
