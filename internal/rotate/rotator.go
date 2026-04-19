// Package rotate provides periodic secret rotation by re-fetching secrets
// from Vault and updating the running process environment.
package rotate

import (
	"context"
	"log/slog"
	"time"
)

// Fetcher is implemented by any type that can return a map of secrets.
type Fetcher interface {
	ReadSecrets(ctx context.Context, paths []string) (map[string]string, error)
}

// Handler is called with fresh secrets after each rotation.
type Handler func(secrets map[string]string)

// Rotator periodically fetches secrets and invokes a Handler.
type Rotator struct {
	fetcher  Fetcher
	paths    []string
	interval time.Duration
	handler  Handler
	logger   *slog.Logger
}

// New creates a Rotator. interval must be positive.
func New(fetcher Fetcher, paths []string, interval time.Duration, handler Handler, logger *slog.Logger) *Rotator {
	if logger == nil {
		logger = slog.Default()
	}
	return &Rotator{
		fetcher:  fetcher,
		paths:    paths,
		interval: interval,
		handler:  handler,
		logger:   logger,
	}
}

// Run starts the rotation loop and blocks until ctx is cancelled.
func (r *Rotator) Run(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			r.logger.Info("rotator stopped")
			return
		case <-ticker.C:
			r.rotate(ctx)
		}
	}
}

func (r *Rotator) rotate(ctx context.Context) {
	secrets, err := r.fetcher.ReadSecrets(ctx, r.paths)
	if err != nil {
		r.logger.Error("rotation fetch failed", "error", err)
		return
	}
	r.logger.Info("secrets rotated", "count", len(secrets))
	r.handler(secrets)
}
