// Package mask provides utilities for redacting secret values from log output.
package mask

import "strings"

// Masker holds a set of secret values and can redact them from strings.
type Masker struct {
	secrets []string
}

// New creates a Masker preloaded with the provided secret values.
// Empty strings are ignored.
func New(secrets map[string]string) *Masker {
	m := &Masker{}
	for _, v := range secrets {
		if v != "" {
			m.secrets = append(m.secrets, v)
		}
	}
	return m
}

// Redact replaces all known secret values in s with "***".
func (m *Masker) Redact(s string) string {
	for _, secret := range m.secrets {
		s = strings.ReplaceAll(s, secret, "***")
	}
	return s
}

// Add registers an additional secret value to be masked.
func (m *Masker) Add(value string) {
	if value != "" {
		m.secrets = append(m.secrets, value)
	}
}

// Len returns the number of tracked secret values.
func (m *Masker) Len() int {
	return len(m.secrets)
}
