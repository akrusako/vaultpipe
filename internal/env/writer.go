// Package env provides utilities for building and exporting environment
// variable sets derived from Vault secrets and OS environment.
package env

import (
	"fmt"
	"strings"
)

// Writer assembles a final environment variable map from multiple sources,
// applying precedence rules: secrets override OS env, explicit overrides win.
type Writer struct {
	base    map[string]string
	overlay map[string]string
}

// New returns a Writer seeded with a copy of the provided base environment
// (typically os.Environ converted to a map).
func New(base map[string]string) *Writer {
	b := make(map[string]string, len(base))
	for k, v := range base {
		b[k] = v
	}
	return &Writer{base: b, overlay: make(map[string]string)}
}

// Apply merges secret key/value pairs into the overlay, overwriting any
// existing overlay entry with the same key.
func (w *Writer) Apply(secrets map[string]string) {
	for k, v := range secrets {
		w.overlay[k] = v
	}
}

// Build returns the merged environment as a slice of "KEY=VALUE" strings
// suitable for exec.Cmd.Env. Overlay values shadow base values.
func (w *Writer) Build() []string {
	merged := make(map[string]string, len(w.base)+len(w.overlay))
	for k, v := range w.base {
		merged[k] = v
	}
	for k, v := range w.overlay {
		merged[k] = v
	}
	out := make([]string, 0, len(merged))
	for k, v := range merged {
		out = append(out, fmt.Sprintf("%s=%s", k, v))
	}
	return out
}

// Lookup returns the effective value for key, checking overlay then base.
func (w *Writer) Lookup(key string) (string, bool) {
	if v, ok := w.overlay[key]; ok {
		return v, true
	}
	v, ok := w.base[key]
	return v, ok
}

// ParseEnvSlice converts a []string of "KEY=VALUE" pairs into a map.
// Entries without '=' are stored with an empty value.
func ParseEnvSlice(env []string) map[string]string {
	m := make(map[string]string, len(env))
	for _, e := range env {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			m[parts[0]] = parts[1]
		} else {
			m[parts[0]] = ""
		}
	}
	return m
}
