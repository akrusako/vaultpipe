// Package metrics provides a simple in-process counter registry used by
// vaultpipe to track operational events across its lifetime.
//
// Counters are safe for concurrent use. A Snapshot can be emitted to a
// structured log or an HTTP debug endpoint at any time.
//
// Well-known counter names:
//
//	"vault.fetch"       – number of Vault secret reads attempted
//	"vault.fetch.error" – number of failed Vault reads
//	"cache.hit"         – secrets served from the local cache
//	"cache.miss"        – secrets that bypassed the cache
//	"exec.start"        – child processes launched
//	"rotate.triggered"  – rotation cycles completed
package metrics
