// Package renew handles Vault token renewal to keep long-running
// processes authenticated throughout their lifetime.
package renew

import (
	"context"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// Renewer periodically renews a Vault token before it expires.
type Renewer struct {
	client   *vaultapi.Client
	interval time.Duration
	logger   interface {
		Infof(format string, args ...interface{})
		Warnf(format string, args ...interface{})
	}
}

// New creates a Renewer that will renew the token at the given interval.
func New(client *vaultapi.Client, interval time.Duration, logger interface {
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
}) *Renewer {
	return &Renewer{
		client:   client,
		interval: interval,
		logger:   logger,
	}
}

// Start begins the renewal loop. It blocks until ctx is cancelled.
func (r *Renewer) Start(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := r.renewOnce(); err != nil {
	 renewal failed: %v", err		r.logger. token renewed successfully")
			}
		}
	}
}

func (r *Renewer) renewOnce() error {
	_, err := r.client.Auth().Token().RenewSelf(0)
	return err
}
