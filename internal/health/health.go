// Package health provides a simple liveness check against a Vault server.
package health

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Status holds the result of a Vault health check.
type Status struct {
	Initialized bool
	Sealed      bool
	Standby     bool
	Code        int
}

// Checker performs health checks against a Vault address.
type Checker struct {
	addr   string
	client *http.Client
}

// New returns a Checker for the given Vault address.
func New(addr string, timeout time.Duration) *Checker {
	return &Checker{
		addr: addr,
		client: &http.Client{Timeout: timeout},
	}
}

// Check calls the Vault /v1/sys/health endpoint and returns a Status.
// Vault returns non-200 codes for standby/sealed states, so we inspect
// the status code directly rather than treating non-200 as an error.
func (c *Checker) Check(ctx context.Context) (*Status, error) {
	url := fmt.Sprintf("%s/v1/sys/health", c.addr)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("health: build request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("health: request failed: %w", err)
	}
	defer resp.Body.Close()

	s := &Status{Code: resp.StatusCode}
	switch resp.StatusCode {
	case http.StatusOK:
		s.Initialized = true
	case 429:
		s.Initialized = true
		s.Standby = true
	case 501:
		// not initialized
	case 503:
		s.Initialized = true
		s.Sealed = true
	default:
		return nil, fmt.Errorf("health: unexpected status %d", resp.StatusCode)
	}
	return s, nil
}

// IsHealthy returns true when Vault is initialized, unsealed, and active.
func (s *Status) IsHealthy() bool {
	return s.Initialized && !s.Sealed && !s.Standby
}
