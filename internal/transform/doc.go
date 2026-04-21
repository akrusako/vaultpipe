// Package transform provides a composable pipeline for transforming secret
// key/value pairs before they are injected into a child process environment.
//
// Usage:
//
//	tr := transform.New(
//		transform.PrefixKeys("APP_"),
//		transform.UppercaseKeys(),
//		transform.TrimValueSpace(),
//	)
//	out, err := tr.Apply(secrets)
//
// Middleware wraps any SecretFetcher so transformations are applied
// transparently at fetch time:
//
//	mw := transform.NewMiddleware(vaultClient.ReadSecrets, tr)
//	env, err := mw.Fetch(paths)
package transform
