// Package coalesce implements priority-ordered secret resolution.
//
// A Coalescer holds an ordered list of secret maps and returns the first
// non-empty value found for a given key. This is useful when secrets may
// originate from multiple Vault paths, environment overrides, or default
// value maps, and the caller wants a single authoritative value without
// writing custom fallback logic at every call site.
//
// Example:
//
//	c := coalesce.New(vaultSecrets, envOverrides, defaults)
//	val, err := c.Lookup("DB_PASSWORD")
//
// Sources are evaluated left-to-right; the first map that contains a
// non-empty string for the key wins. Use Merge to collapse all sources
// into a single flat map for injection into a child process environment.
package coalesce
