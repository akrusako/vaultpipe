// Package scrub provides middleware for stripping sensitive keys from secret
// maps before they are injected into a child process environment.
package scrub

import "strings"

// Scrubber removes secret entries whose keys match a deny-list of patterns.
type Scrubber struct {
	patterns []string
}

// Option configures a Scrubber.
type Option func(*Scrubber)

// WithPatterns appends additional deny-list patterns (case-insensitive prefix
// or exact match).
func WithPatterns(patterns ...string) Option {
	return func(s *Scrubber) {
		for _, p := range patterns {
			s.patterns = append(s.patterns, strings.ToLower(p))
		}
	}
}

// defaultPatterns are always denied regardless of user configuration.
var defaultPatterns = []string{
	"vault_token",
	"vault_addr",
	"vaultpipe_",
}

// New creates a Scrubber with the built-in deny-list plus any extra patterns
// supplied via options.
func New(opts ...Option) *Scrubber {
	s := &Scrubber{}
	for _, p := range defaultPatterns {
		s.patterns = append(s.patterns, strings.ToLower(p))
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Apply returns a copy of secrets with denied keys removed.
func (s *Scrubber) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if !s.denied(k) {
			out[k] = v
		}
	}
	return out
}

// denied reports whether key matches any deny-list pattern.
func (s *Scrubber) denied(key string) bool {
	lower := strings.ToLower(key)
	for _, p := range s.patterns {
		if lower == p || strings.HasPrefix(lower, p) {
			return true
		}
	}
	return false
}
