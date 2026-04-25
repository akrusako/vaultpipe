// Package validate provides configurable validation of secret values fetched
// from Vault before they are injected into a child process environment.
//
// # Usage
//
//	v := validate.New(
//		validate.NonEmpty(),
//		validate.MaxLength(4096),
//		validate.MatchesPattern("jwt", `^[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+$`),
//	)
//
//	mw := validate.NewMiddleware(v, myFetcher)
//	secrets, err := mw(paths)
//
// Rules are evaluated in registration order. The first failing rule for a
// given key short-circuits further evaluation for that key. All failing keys
// are collected and returned as a single error from Check or the middleware.
package validate
