package validate

import "fmt"

// SecretFetcher is the function signature used to retrieve secrets from Vault.
type SecretFetcher func(paths []string) (map[string]string, error)

// Middleware wraps a SecretFetcher and validates the returned secrets before
// passing them to the caller. If any secret fails validation the fetch is
// aborted and a combined error is returned.
func NewMiddleware(v *Validator, next SecretFetcher) SecretFetcher {
	return func(paths []string) (map[string]string, error) {
		secrets, err := next(paths)
		if err != nil {
			return nil, err
		}
		if errs := v.Check(secrets); errs != nil {
			return nil, fmt.Errorf("secret validation failed: %w", buildError(errs))
		}
		return secrets, nil
	}
}

// buildError collapses the per-key error map into a single descriptive error.
func buildError(errs map[string]error) error {
	var msg string
	for k, e := range errs {
		if msg != "" {
			msg += "; "
		}
		msg += k + ": " + e.Error()
	}
	return fmt.Errorf("%s", msg)
}
