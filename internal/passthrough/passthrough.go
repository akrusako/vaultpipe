// Package passthrough selectively forwards environment variables from the
// host environment into the child process, applying an optional allowlist
// or denylist of key patterns.
package passthrough

import (
	"os"
	"strings"
)

// Option configures a Passthrough.
type Option func(*Passthrough)

// Passthrough filters host environment variables before they reach the child
// process.
type Passthrough struct {
	allow []string
	deny  []string
}

// WithAllow restricts forwarding to keys whose names contain one of the given
// substrings (case-insensitive).
func WithAllow(patterns ...string) Option {
	return func(p *Passthrough) {
		p.allow = append(p.allow, patterns...)
	}
}

// WithDeny blocks keys whose names contain one of the given substrings
// (case-insensitive).
func WithDeny(patterns ...string) Option {
	return func(p *Passthrough) {
		p.deny = append(p.deny, patterns...)
	}
}

// New returns a Passthrough configured with the provided options.
func New(opts ...Option) *Passthrough {
	p := &Passthrough{}
	for _, o := range opts {
		o(p)
	}
	return p
}

// Filter returns a copy of env containing only the entries that pass the
// allow/deny rules. If no allow patterns are set every key is considered
// allowed before the deny check is applied.
func (p *Passthrough) Filter(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if !p.allowed(k) {
			continue
		}
		if p.denied(k) {
			continue
		}
		out[k] = v
	}
	return out
}

// FromOS returns the current process environment as a map, then applies
// Filter to it.
func (p *Passthrough) FromOS() map[string]string {
	raw := os.Environ()
	env := make(map[string]string, len(raw))
	for _, pair := range raw {
		if idx := strings.IndexByte(pair, '='); idx >= 0 {
			env[pair[:idx]] = pair[idx+1:]
		}
	}
	return p.Filter(env)
}

func (p *Passthrough) allowed(key string) bool {
	if len(p.allow) == 0 {
		return true
	}
	lower := strings.ToLower(key)
	for _, pat := range p.allow {
		if strings.Contains(lower, strings.ToLower(pat)) {
			return true
		}
	}
	return false
}

func (p *Passthrough) denied(key string) bool {
	lower := strings.ToLower(key)
	for _, pat := range p.deny {
		if strings.Contains(lower, strings.ToLower(pat)) {
			return true
		}
	}
	return false
}
