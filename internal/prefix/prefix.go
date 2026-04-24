// Package prefix provides utilities for adding a namespace prefix to
// secret keys before they are injected into a child process environment.
package prefix

import "strings"

// Prefixer adds a fixed string prefix to every key in a secrets map.
type Prefixer struct {
	prefix string
}

// Option is a functional option for Prefixer.
type Option func(*Prefixer)

// WithPrefix sets the prefix string. Leading and trailing whitespace is
// trimmed; an empty string is a no-op prefix.
func WithPrefix(p string) Option {
	return func(pr *Prefixer) {
		pr.prefix = strings.TrimSpace(p)
	}
}

// New returns a new Prefixer configured by opts.
func New(opts ...Option) *Prefixer {
	pr := &Prefixer{}
	for _, o := range opts {
		o(pr)
	}
	return pr
}

// Apply returns a new map whose keys are the original keys prepended with
// the configured prefix. Values are copied without modification. If the
// prefix is empty the original map is returned unchanged.
func (p *Prefixer) Apply(secrets map[string]string) map[string]string {
	if p.prefix == "" {
		return secrets
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[p.prefix+k] = v
	}
	return out
}

// Strip returns a new map with the prefix removed from every key that
// carries it. Keys that do not start with the prefix are dropped.
func (p *Prefixer) Strip(secrets map[string]string) map[string]string {
	if p.prefix == "" {
		out := make(map[string]string, len(secrets))
		for k, v := range secrets {
			out[k] = v
		}
		return out
	}
	out := make(map[string]string)
	for k, v := range secrets {
		if strings.HasPrefix(k, p.prefix) {
			out[strings.TrimPrefix(k, p.prefix)] = v
		}
	}
	return out
}
