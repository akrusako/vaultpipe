// Package throttle implements a semaphore-based concurrency limiter for
// outbound Vault secret-fetch operations.
//
// # Overview
//
// When vaultpipe resolves a large number of secret paths in parallel it can
// inadvertently flood the Vault server with simultaneous requests. The Limiter
// type caps the number of in-flight fetches to a configurable maximum, causing
// excess callers to block until a slot becomes available or their context
// expires.
//
// # Usage
//
//	cfg := throttle.DefaultConfig()   // Concurrency: 8
//	limiter, err := throttle.New(cfg)
//
//	// Wrap an existing fetch function:
//	mw, err := throttle.NewMiddleware(limiter, myFetchFunc)
//	secrets, err := mw.Fetch(ctx, paths)
package throttle
