// Package debounce provides a debouncing wrapper around a secret Fetcher.
//
// When multiple vault-watch events fire within a short window (e.g. a bulk
// secret rotation), debounce coalesces the resulting fetch calls so that the
// underlying Vault API is queried only once per quiet period. This reduces
// load on Vault and prevents redundant process restarts.
//
// Basic usage:
//
//	d, err := debounce.New(200*time.Millisecond, myFetcher)
//	if err != nil {
//		log.Fatal(err)
//	}
//	data, err := d.Fetch(ctx, paths)
package debounce
