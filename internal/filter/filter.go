// Package filter provides key-based filtering for secrets before they are
// injected into a process environment.
package filter

import "strings"

// Filter holds include and exclude glob-style prefix rules for secret keys.
type Filter struct {
	include []string
	exclude []string
}

// New creates a Filter with the given include and exclude prefix lists.
// If include is empty, all keys are included by default.
func New(include, exclude []string) *Filter {
	return &Filter{
		include: include,
		exclude: exclude,
	}
}

// Allow returns true if the given key should be passed through.
func (f *Filter) Allow(key string) bool {
	if !f.matchesInclude(key) {
		return false
	}
	for _, ex := range f.exclude {
		if matchPrefix(key, ex) {
			return false
		}
	}
	return true
}

// Apply filters a map of secrets, returning only allowed keys.
func (f *Filter) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if f.Allow(k) {
			out[k] = v
		}
	}
	return out
}

func (f *Filter) matchesInclude(key string) bool {
	if len(f.include) == 0 {
		return true
	}
	for _, in := range f.include {
		if matchPrefix(key, in) {
			return true
		}
	}
	return false
}

func matchPrefix(key, pattern string) bool {
	return strings.HasPrefix(strings.ToUpper(key), strings.ToUpper(pattern))
}
