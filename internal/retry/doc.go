// Package retry implements exponential backoff retry logic for use
// throughout vaultpipe when communicating with Vault or other external
// systems that may experience transient failures.
//
// Usage:
//
//	err := retry.Do(ctx, retry.DefaultConfig(), func() error {
//		return doSomething()
//	})
//
// To prevent retrying on a known-fatal error, wrap it with retry.Permanent:
//
//	return retry.Permanent(err)
package retry
