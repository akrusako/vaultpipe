package watch

import (
	"errors"
	"time"
)

// Config holds tunable parameters for a Watcher.
type Config struct {
	// Interval is how often Vault is polled for secret changes.
	Interval time.Duration

	// Paths is the list of Vault secret paths to watch.
	Paths []string
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Interval: 30 * time.Second,
	}
}

// Validate returns an error if the Config is not usable.
func (c Config) Validate() error {
	if c.Interval <= 0 {
		return errors.New("watch: interval must be positive")
	}
	if len(c.Paths) == 0 {
		return errors.New("watch: at least one secret path is required")
	}
	return nil
}

// NewFromConfig constructs a Watcher from a Config.
func NewFromConfig(cfg Config, fetcher SecretFetcher, onChange ChangeHandler) (*Watcher, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return New(fetcher, cfg.Paths, cfg.Interval, onChange), nil
}
