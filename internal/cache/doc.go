// Package cache implements a thread-safe, TTL-based in-memory cache for
// Vault secret paths. It is used by the Vault client to avoid redundant
// network requests when the same secret path is read multiple times within
// a short window.
//
// Entries are automatically considered stale once their TTL has elapsed;
// callers should re-fetch from Vault and repopulate the cache on a miss.
package cache
