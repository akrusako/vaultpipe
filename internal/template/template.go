// Package template provides secret path template rendering,
// allowing users to specify dynamic paths using environment variables.
package template

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

// Renderer renders secret path templates using environment variables
// and a provided data map as the template context.
type Renderer struct {
	env map[string]string
}

// New creates a new Renderer seeded with the current OS environment.
func New() *Renderer {
	env := make(map[string]string)
	for _, e := range os.Environ() {
		for i := 0; i < len(e); i++ {
			if e[i] == '=' {
				env[e[:i]] = e[i+1:]
				break
			}
		}
	}
	return &Renderer{env: env}
}

// Render executes the given path template using the renderer's environment
// map as the data context. Template variables use Go text/template syntax,
// e.g. "secret/{{.ENVIRONMENT}}/db".
func (r *Renderer) Render(path string) (string, error) {
	tmpl, err := template.New("path").Option("missingkey=error").Parse(path)
	if err != nil {
		return "", fmt.Errorf("template parse error for %q: %w", path, err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, r.env); err != nil {
		return "", fmt.Errorf("template render error for %q: %w", path, err)
	}
	return buf.String(), nil
}

// RenderAll renders a slice of path templates, returning the rendered paths
// or the first error encountered.
func (r *Renderer) RenderAll(paths []string) ([]string, error) {
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		rendered, err := r.Render(p)
		if err != nil {
			return nil, err
		}
		out = append(out, rendered)
	}
	return out, nil
}
