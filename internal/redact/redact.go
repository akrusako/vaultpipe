// Package redact provides utilities for scrubbing sensitive values
// from log output and error messages before they are written anywhere.
package redact

import "strings"

// Redactor replaces known secret values with a placeholder string.
type Redactor struct {
	secrets     []string
	placeholder string
}

// New returns a Redactor that replaces each value in secrets with
// placeholder. Empty secret values are silently ignored.
func New(placeholder string, secrets []string) *Redactor {
	filtered := make([]string, 0, len(secrets))
	for _, s := range secrets {
		if strings.TrimSpace(s) != "" {
			filtered = append(filtered, s)
		}
	}
	return &Redactor{secrets: filtered, placeholder: placeholder}
}

// Scrub returns a copy of s with all registered secret values replaced
// by the placeholder.
func (r *Redactor) Scrub(s string) string {
	for _, secret := range r.secrets {
		s = strings.ReplaceAll(s, secret, r.placeholder)
	}
	return s
}

// Add registers an additional secret value. Empty values are ignored.
func (r *Redactor) Add(secret string) {
	if strings.TrimSpace(secret) != "" {
		r.secrets = append(r.secrets, secret)
	}
}

// Len returns the number of registered secrets.
func (r *Redactor) Len() int {
	return len(r.secrets)
}
