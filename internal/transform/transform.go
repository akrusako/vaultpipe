// Package transform provides key/value transformation pipelines
// for secret data before injection into process environments.
package transform

import (
	"fmt"
	"strings"
)

// Func is a transformation applied to a secret key or value.
type Func func(key, value string) (string, string, error)

// Transformer applies an ordered chain of Funcs to secret maps.
type Transformer struct {
	funcs []Func
}

// New returns a Transformer that applies the given Funcs in order.
func New(fns ...Func) *Transformer {
	return &Transformer{funcs: fns}
}

// Apply runs every registered Func over each key/value pair in secrets.
// If any Func returns an error the call is aborted and the error returned.
func (t *Transformer) Apply(secrets map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		curKey, curVal := k, v
		for _, fn := range t.funcs {
			nk, nv, err := fn(curKey, curVal)
			if err != nil {
				return nil, fmt.Errorf("transform %q: %w", curKey, err)
			}
			curKey, curVal = nk, nv
		}
		out[curKey] = curVal
	}
	return out, nil
}

// PrefixKeys returns a Func that prepends prefix to every key.
func PrefixKeys(prefix string) Func {
	return func(key, value string) (string, string, error) {
		return prefix + key, value, nil
	}
}

// UppercaseKeys returns a Func that upper-cases every key.
func UppercaseKeys() Func {
	return func(key, value string) (string, string, error) {
		return strings.ToUpper(key), value, nil
	}
}

// TrimValueSpace returns a Func that trims leading/trailing whitespace from values.
func TrimValueSpace() Func {
	return func(key, value string) (string, string, error) {
		return key, strings.TrimSpace(value), nil
	}
}
