// Package expire implements TTL-based expiration tracking for secrets
// fetched from Vault.
//
// # Overview
//
// A Tracker stores named entries, each associated with a string value and
// an absolute expiration deadline. Callers can set, retrieve, delete, and
// purge entries. Expired entries are invisible to Get but remain in memory
// until Purge is called.
//
// # Middleware
//
// NewMiddleware wraps a Fetcher and transparently caches its results in a
// Tracker for a configurable TTL. This prevents redundant round-trips to
// Vault when the same secret paths are requested within the expiry window.
//
// Example:
//
//	tracker := expire.New(nil)
//	mw := expire.NewMiddleware(vaultFetch, tracker, 30*time.Second)
//	secrets, err := mw.Fetch(ctx, []string{"secret/app/db"})
package expire
