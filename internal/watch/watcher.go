// Package watch provides periodic secret refresh by re-fetching secrets
// from Vault on a configurable interval and notifying callers of changes.
package watch

import (
	"context"
	"time"

	"github.com/your-org/vaultpipe/internal/diff"
	"github.com/your-org/vaultpipe/internal/snapshot"
)

// SecretFetcher is the interface for fetching secrets from Vault.
type SecretFetcher interface {
	ReadSecrets(ctx context.Context, paths []string) (map[string]string, error)
}

// ChangeHandler is called whenever a change in secrets is detected.
type ChangeHandler func(added, removed, changed map[string]string)

// Watcher polls Vault for secret changes on a fixed interval.
type Watcher struct {
	fetcher  SecretFetcher
	paths    []string
	interval time.Duration
	snap     *snapshot.Snapshot
	diff     *diff.Differ
	onChange ChangeHandler
}

// New creates a Watcher that polls the given paths every interval,
// calling onChange when secrets differ from the previous fetch.
func New(fetcher SecretFetcher, paths []string, interval time.Duration, onChange ChangeHandler) *Watcher {
	return &Watcher{
		fetcher:  fetcher,
		paths:    paths,
		interval: interval,
		snap:     snapshot.New(),
		diff:     diff.New(),
		onChange: onChange,
	}
}

// Watch starts the polling loop and blocks until ctx is cancelled.
func (w *Watcher) Watch(ctx context.Context) error {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := w.poll(ctx); err != nil {
				// Non-fatal: log elsewhere; keep watching.
				continue
			}
		}
	}
}

func (w *Watcher) poll(ctx context.Context) error {
	secrets, err := w.fetcher.ReadSecrets(ctx, w.paths)
	if err != nil {
		return err
	}

	if !w.snap.Changed(secrets) {
		return nil
	}

	prev := w.snap.Latest()
	w.snap.Save(secrets)

	result := w.diff.Compare(prev, secrets)
	if result.HasChanges() && w.onChange != nil {
		w.onChange(result.Added, result.Removed, result.Changed)
	}
	return nil
}
