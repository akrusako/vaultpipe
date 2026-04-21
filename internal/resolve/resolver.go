// Package resolve provides utilities for resolving secret path templates
// and expanding environment variable references within Vault secret paths.
package resolve

import (
	"fmt"
	"os"
	"strings"
)

// Resolver expands variable references in secret path strings.
// References use the ${VAR} or $VAR syntax and are resolved
// against a provided environment map, falling back to os.Getenv.
type Resolver struct {
	env map[string]string
}

// New returns a Resolver that consults env first, then the process
// environment for any variable not found in env.
func New(env map[string]string) *Resolver {
	copy := make(map[string]string, len(env))
	for k, v := range env {
		copy[k] = v
	}
	return &Resolver{env: copy}
}

// Expand resolves all variable references within path and returns the
// expanded string. An error is returned if a referenced variable is
// not defined in either the resolver's env map or the OS environment.
func (r *Resolver) Expand(path string) (string, error) {
	var expandErr error
	result := os.Expand(path, func(key string) string {
		if expandErr != nil {
			return ""
		}
		if v, ok := r.env[key]; ok {
			return v
		}
		if v, ok := os.LookupEnv(key); ok {
			return v
		}
		expandErr = fmt.Errorf("resolve: undefined variable %q in path %q", key, path)
		return ""
	})
	if expandErr != nil {
		return "", expandErr
	}
	return result, nil
}

// ExpandAll resolves variable references in each path and returns the
// expanded slice. The first unresolvable variable causes an early return.
func (r *Resolver) ExpandAll(paths []string) ([]string, error) {
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		expanded, err := r.Expand(p)
		if err != nil {
			return nil, err
		}
		out = append(out, expanded)
	}
	return out, nil
}

// HasVars reports whether path contains at least one variable reference.
func HasVars(path string) bool {
	return strings.Contains(path, "$")
}
