package drain

import (
	"context"
	"fmt"
)

// Fetcher is any function that retrieves secrets from a path.
type Fetcher func(ctx context.Context, paths []string) (map[string]string, error)

// NewMiddleware wraps a Fetcher so that every call is tracked by the
// provided Drainer. If the Drainer is already closed the call is
// rejected immediately with an error.
func NewMiddleware(d *Drainer, next Fetcher) Fetcher {
	return func(ctx context.Context, paths []string) (map[string]string, error) {
		if !d.Acquire() {
			return nil, fmt.Errorf("drain: drainer is closed, refusing new fetch")
		}
		defer d.Release()
		return next(ctx, paths)
	}
}
