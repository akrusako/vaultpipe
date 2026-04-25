// Package dedupe implements key-value deduplication for secret maps.
//
// It is used in watch/rotation pipelines to avoid re-injecting or
// re-notifying on secrets that have not changed between fetch cycles.
//
// Basic usage:
//
//	d := dedupe.New()
//
//	// First fetch — all keys are new.
//	changed := d.Filter(fetchedSecrets)
//
//	// Subsequent fetch — only modified or new keys are returned.
//	changed = d.Filter(fetchedSecrets)
//
// Call Reset to force all keys to be treated as new on the next Filter call.
package dedupe
