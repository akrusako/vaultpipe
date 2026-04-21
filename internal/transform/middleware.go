package transform

// SecretFetcher is a function that retrieves secrets by path.
type SecretFetcher func(paths []string) (map[string]string, error)

// Middleware wraps a SecretFetcher and applies a Transformer to every
// result before returning it to the caller.
type Middleware struct {
	fetcher     SecretFetcher
	transformer *Transformer
}

// NewMiddleware returns a Middleware that applies t to secrets returned by f.
func NewMiddleware(f SecretFetcher, t *Transformer) *Middleware {
	return &Middleware{fetcher: f, transformer: t}
}

// Fetch retrieves secrets via the wrapped fetcher and transforms them.
func (m *Middleware) Fetch(paths []string) (map[string]string, error) {
	secrets, err := m.fetcher(paths)
	if err != nil {
		return nil, err
	}
	return m.transformer.Apply(secrets)
}
