// Package ratelimit implements a token-bucket rate limiter used to throttle
// outbound requests to HashiCorp Vault.
//
// The limiter is safe for concurrent use. Callers invoke Wait before each
// Vault operation; Wait blocks until a token is available or the supplied
// context is cancelled.
//
// Example:
//
//	lim, err := ratelimit.New(ratelimit.DefaultConfig())
//	if err != nil {
//		log.Fatal(err)
//	}
//	if err := lim.Wait(ctx); err != nil {
//		return err
//	}
//	// perform Vault read …
package ratelimit
