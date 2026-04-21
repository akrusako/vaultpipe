// Package quota provides a sliding-window read-rate limiter for Vault secret
// paths. It is designed to prevent runaway polling loops or misconfigured
// watchers from hammering the Vault API.
//
// Basic usage:
//
//	cfg := quota.DefaultConfig()
//	cfg.MaxReads = 60
//	cfg.Window  = time.Minute
//
//	q, err := quota.New(cfg)
//	if err != nil { ... }
//
//	// Wrap an existing fetch function:
//	guarded := quota.NewMiddleware(q, myReadFunc)
//
// The quota is tracked per path and resets automatically as timestamps slide
// out of the configured window. Call Reset to clear a path explicitly.
package quota
