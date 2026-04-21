// Package sanitize provides utilities for cleaning and validating
// environment variable keys before they are injected into a child process.
package sanitize

import (
	"fmt"
	"regexp"
	"strings"
)

var validKey = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// Sanitizer cleans and validates environment variable keys.
type Sanitizer struct {
	upper bool
}

// Option configures a Sanitizer.
type Option func(*Sanitizer)

// WithUppercase returns an Option that uppercases all keys.
func WithUppercase() Option {
	return func(s *Sanitizer) { s.upper = true }
}

// New creates a new Sanitizer with the given options.
func New(opts ...Option) *Sanitizer {
	s := &Sanitizer{}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Key normalises a single environment variable key. It replaces hyphens and
// dots with underscores, optionally uppercases the result, and returns an
// error if the resulting key is not a valid POSIX identifier.
func (s *Sanitizer) Key(raw string) (string, error) {
	k := strings.NewReplacer("-", "_", ".", "_", "/", "_").Replace(raw)
	if s.upper {
		k = strings.ToUpper(k)
	}
	if !validKey.MatchString(k) {
		return "", fmt.Errorf("sanitize: %q is not a valid environment variable key", k)
	}
	return k, nil
}

// Map sanitises all keys in the supplied map, returning a new map with clean
// keys. Values are left unchanged. Keys that fail validation are skipped and
// their original key is collected into the returned error slice.
func (s *Sanitizer) Map(in map[string]string) (map[string]string, []error) {
	out := make(map[string]string, len(in))
	var errs []error
	for k, v := range in {
		clean, err := s.Key(k)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		out[clean] = v
	}
	return out, errs
}
