// Package diff provides utilities for detecting changes between two snapshots
// of secret key-value maps. It is used by the rotation pipeline to determine
// whether a process restart or environment refresh is necessary after a
// Vault secret lease renewal.
package diff
