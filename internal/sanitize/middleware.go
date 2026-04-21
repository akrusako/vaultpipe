package sanitize

// Middleware wraps a downstream secret consumer, sanitising keys before
// passing them on. It is designed to sit between the Vault client and the
// env.Writer so that normalisation is applied transparently.
type Middleware struct {
	s    *Sanitizer
	next func(map[string]string) error
}

// NewMiddleware creates a Middleware that sanitises keys with s and forwards
// the cleaned map to next.
func NewMiddleware(s *Sanitizer, next func(map[string]string) error) *Middleware {
	return &Middleware{s: s, next: next}
}

// Apply sanitises the keys in secrets and calls the downstream handler.
// Any keys that fail validation are dropped; the first sanitisation error is
// returned alongside any error from the downstream handler.
func (m *Middleware) Apply(secrets map[string]string) error {
	clean, errs := m.s.Map(secrets)
	var firstErr error
	if len(errs) > 0 {
		firstErr = errs[0]
	}
	if err := m.next(clean); err != nil {
		return err
	}
	return firstErr
}
