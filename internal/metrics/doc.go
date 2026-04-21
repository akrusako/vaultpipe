// Package metrics provides a simple in-process counter registry used by
// vaultpipe to track operational events across its lifetime.
//
// Counters are safe for concurrent use. A Snapshot can be emitted to a
// structured log or an HTTP debug endpoint at any time.
//
// Well-known counter names:
//
//	"vault.fetch"        – number of Vault secret reads attempted
//	"vault.fetch.error"  – number of failed Vault reads
//	"cache.hit"          – secrets served from the local cache
//	"cache.miss"         – secrets that bypassed the cache
//	"cache.evict"        – entries removed from the cache (TTL or capacity)
//	"exec.start"         – child processes launched
//	"exec.error"         – child processes that exited with a non-zero status
//	"rotate.triggered"   – rotation cycles initiated
//	"rotate.error"       – rotation cycles that failed to complete
//
// Counter names follow a "component.event" convention. Errors are always
// represented as a sub-key of the operation they belong to (e.g.
// "vault.fetch.error" is the error counterpart of "vault.fetch").
package metrics
