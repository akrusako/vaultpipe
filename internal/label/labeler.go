// Package label attaches a fixed set of key/value metadata labels to every
// secret map produced by a fetcher. Labels are applied after the fetch and
// before any downstream middleware sees the data, making them useful for
// environment tagging (e.g. app=myapp, env=production).
package label

import (
	"context"
	"fmt"
	"strings"
)

// Fetcher is the interface satisfied by any secret source.
type Fetcher interface {
	Fetch(ctx context.Context, paths []string) (map[string]string, error)
}

// Labeler wraps a Fetcher and merges a static label map into every result.
type Labeler struct {
	next   Fetcher
	labels map[string]string
}

// Option configures a Labeler.
type Option func(*Labeler)

// WithLabel adds a single label. The key is trimmed of whitespace. An empty
// key or value is silently ignored.
func WithLabel(key, value string) Option {
	return func(l *Labeler) {
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key != "" && value != "" {
			l.labels[key] = value
		}
	}
}

// New returns a Labeler that decorates next with the provided labels.
// Labels must not be nil; pass no Option to create a pass-through wrapper.
func New(next Fetcher, opts ...Option) (*Labeler, error) {
	if next == nil {
		return nil, fmt.Errorf("label: next fetcher must not be nil")
	}
	l := &Labeler{
		next:   next,
		labels: make(map[string]string),
	}
	for _, o := range opts {
		o(l)
	}
	return l, nil
}

// Fetch delegates to the wrapped Fetcher and then merges static labels into
// the returned map. Labels never overwrite secrets that share the same key;
// the secret value always wins.
func (l *Labeler) Fetch(ctx context.Context, paths []string) (map[string]string, error) {
	result, err := l.next.Fetch(ctx, paths)
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(result)+len(l.labels))
	for k, v := range l.labels {
		out[k] = v
	}
	// Secrets overwrite labels when keys collide.
	for k, v := range result {
		out[k] = v
	}
	return out, nil
}

// Labels returns a copy of the currently registered label map.
func (l *Labeler) Labels() map[string]string {
	copy := make(map[string]string, len(l.labels))
	for k, v := range l.labels {
		copy[k] = v
	}
	return copy
}
