// Package flatten provides utilities for flattening nested secret maps
// into a single-level key=value map suitable for environment injection.
package flatten

import (
	"fmt"
	"strings"
)

// Options controls how keys are constructed during flattening.
type Options struct {
	// Separator is placed between nested key segments. Defaults to "_".
	Separator string
	// Uppercase converts all keys to uppercase when true.
	Uppercase bool
}

// DefaultOptions returns a sensible default Options value.
func DefaultOptions() Options {
	return Options{
		Separator: "_",
		Uppercase: true,
	}
}

// Flattener collapses nested map[string]any structures into a flat
// map[string]string.
type Flattener struct {
	opts Options
}

// New creates a Flattener with the given options.
func New(opts Options) *Flattener {
	if opts.Separator == "" {
		opts.Separator = "_"
	}
	return &Flattener{opts: opts}
}

// Flatten converts a potentially nested map into a flat map[string]string.
// Nested keys are joined with the configured separator.
func (f *Flattener) Flatten(input map[string]any) map[string]string {
	out := make(map[string]string)
	f.flatten("", input, out)
	return out
}

func (f *Flattener) flatten(prefix string, input map[string]any, out map[string]string) {
	for k, v := range input {
		fullKey := f.buildKey(prefix, k)
		switch val := v.(type) {
		case map[string]any:
			f.flatten(fullKey, val, out)
		case map[string]string:
			for ik, iv := range val {
				out[f.buildKey(fullKey, ik)] = iv
			}
		default:
			out[fullKey] = fmt.Sprintf("%v", val)
		}
	}
}

func (f *Flattener) buildKey(prefix, key string) string {
	var k string
	if prefix == "" {
		k = key
	} else {
		k = prefix + f.opts.Separator + key
	}
	if f.opts.Uppercase {
		return strings.ToUpper(k)
	}
	return k
}
