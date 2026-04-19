// Package config loads and validates vaultpipe runtime configuration
// from environment variables.
package config

import (
	"errors"
	"os"
	"strings"
	"time"
)

// Config holds all runtime configuration for vaultpipe.
type Config struct {
	VaultAddr    string
	VaultToken   string
	SecretPaths  []string
	Command      []string
	RenewInterval time.Duration
}

// FromEnv builds a Config from environment variables.
//
//	VAULT_ADDR        – Vault server address (required)
//	VAULT_TOKEN       – Vault token (required)
//	VAULTPIPE_SECRETS – comma-separated secret paths (required)
//	VAULTPIPE_RENEW   – token renewal interval (default: 5m)
func FromEnv() (*Config, error) {
	cfg := &Config{
		VaultAddr:  os.Getenv("VAULT_ADDR"),
		VaultToken: os.Getenv("VAULT_TOKEN"),
	}

	if raw := os.Getenv("VAULTPIPE_SECRETS"); raw != "" {
		for _, p := range strings.Split(raw, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				cfg.SecretPaths = append(cfg.SecretPaths, p)
			}
		}
	}

	interval := os.Getenv("VAULTPIPE_RENEW")
	if interval == "" {
		interval = "5m"
	}
	var err error
	cfg.RenewInterval, err = time.ParseDuration(interval)
	if err != nil {
		return nil, errors.New("VAULTPIPE_RENEW: invalid duration: " + interval)
	}

	return cfg, nil
}

// Validate returns an error if any required field is missing.
func (c *Config) Validate() error {
	if c.VaultAddr == "" {
		return errors.New("VAULT_ADDR is required")
	}
	if c.VaultToken == "" {
		return errors.New("VAULT_TOKEN is required")
	}
	if len(c.SecretPaths) == 0 {
		return errors.New("VAULTPIPE_SECRETS is required")
	}
	if len(c.Command) == 0 {
		return errors.New("command is required")
	}
	return nil
}
