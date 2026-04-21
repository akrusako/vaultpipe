package quota

import "fmt"

// ReadFunc is the signature of a secret-fetching function.
type ReadFunc func(path string) (map[string]string, error)

// NewMiddleware wraps a ReadFunc with quota enforcement.
// Each call to the returned function checks the quota for the given path
// before delegating to next. If the quota is exceeded the underlying
// function is not called and ErrQuotaExceeded is returned.
func NewMiddleware(q *Quota, next ReadFunc) ReadFunc {
	return func(path string) (map[string]string, error) {
		if err := q.Check(path); err != nil {
			return nil, fmt.Errorf("quota middleware: %w", err)
		}
		return next(path)
	}
}
