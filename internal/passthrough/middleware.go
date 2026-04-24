package passthrough

// SecretFetcher is a function that retrieves secrets from Vault and returns
// them as a key/value map.
type SecretFetcher func() (map[string]string, error)

// Middleware wraps a SecretFetcher and merges the filtered host environment
// into the returned secrets map. Host variables are added only when no secret
// key already exists with the same name, preserving Vault values as the
// authoritative source.
type Middleware struct {
	pt   *Passthrough
	next SecretFetcher
}

// NewMiddleware returns a Middleware that applies pt to the host environment
// and merges the result into secrets returned by next.
func NewMiddleware(pt *Passthrough, next SecretFetcher) *Middleware {
	if pt == nil {
		pt = New()
	}
	return &Middleware{pt: pt, next: next}
}

// Fetch calls the underlying SecretFetcher, then injects allowed host
// environment variables that are not already present in the secrets map.
func (m *Middleware) Fetch() (map[string]string, error) {
	secrets, err := m.next()
	if err != nil {
		return nil, err
	}

	host := m.pt.FromOS()
	merged := make(map[string]string, len(secrets)+len(host))

	// Host env goes in first so Vault secrets can overwrite.
	for k, v := range host {
		merged[k] = v
	}
	for k, v := range secrets {
		merged[k] = v
	}
	return merged, nil
}
