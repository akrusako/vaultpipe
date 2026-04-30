// Package pin provides a mechanism to lock a secret map to a known-good
// snapshot and reject any fetch result that diverges beyond a configured
// threshold. This guards against accidental or malicious secret rotation
// that would silently change the environment seen by a child process.
package pin

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// ErrPinViolation is returned when a fetched secret map differs from the
// pinned baseline by more keys than MaxDrift allows.
var ErrPinViolation = errors.New("pin: secret map violates pinned baseline")

// Fetcher is the interface satisfied by any secret source.
type Fetcher interface {
	Fetch(ctx context.Context, paths []string) (map[string]string, error)
}

// Config controls pin behaviour.
type Config struct {
	// MaxDrift is the maximum number of keys that may differ between the
	// pinned baseline and a freshly fetched map before ErrPinViolation is
	// returned. Zero means no drift is tolerated.
	MaxDrift int
}

// Pinner wraps a Fetcher and enforces a pinned baseline.
type Pinner struct {
	cfg     Config
	inner   Fetcher
	mu      sync.Mutex
	baseline map[string]string
}

// New returns a Pinner that delegates to inner and enforces cfg.
func New(inner Fetcher, cfg Config) (*Pinner, error) {
	if inner == nil {
		return nil, errors.New("pin: inner fetcher must not be nil")
	}
	if cfg.MaxDrift < 0 {
		return nil, errors.New("pin: MaxDrift must be >= 0")
	}
	return &Pinner{cfg: cfg, inner: inner}, nil
}

// Fetch retrieves secrets via the inner fetcher. On the first call the result
// is stored as the baseline. Subsequent calls compare against the baseline and
// return ErrPinViolation when drift exceeds MaxDrift.
func (p *Pinner) Fetch(ctx context.Context, paths []string) (map[string]string, error) {
	result, err := p.inner.Fetch(ctx, paths)
	if err != nil {
		return nil, err
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.baseline == nil {
		p.baseline = copyMap(result)
		return result, nil
	}

	drift := countDrift(p.baseline, result)
	if drift > p.cfg.MaxDrift {
		return nil, fmt.Errorf("%w: %d key(s) changed (max %d)", ErrPinViolation, drift, p.cfg.MaxDrift)
	}
	return result, nil
}

// Reset clears the pinned baseline so the next Fetch call establishes a new one.
func (p *Pinner) Reset() {
	p.mu.Lock()
	p.baseline = nil
	p.mu.Unlock()
}

func countDrift(baseline, current map[string]string) int {
	seen := make(map[string]struct{}, len(baseline))
	drift := 0
	for k, v := range current {
		seen[k] = struct{}{}
		if bv, ok := baseline[k]; !ok || bv != v {
			drift++
		}
	}
	for k := range baseline {
		if _, ok := seen[k]; !ok {
			drift++
		}
	}
	return drift
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
