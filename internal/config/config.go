// Package config loads and validates vaultpipe configuration from environment
// variables. All configuration is intentionally sourced from the environment
// to avoid secrets appearing in config files or command-line arguments.
package config

import (
	"errors"
	"os"
	"strings"
)

// Config holds all runtime configuration for vaultpipe.
type Config struct {
	// VaultAddr is the address of the Vault server.
	VaultAddr string
	// VaultToken is the token used to authenticate with Vault.
	VaultToken string
	// SecretPaths is the list of Vault secret paths (or path templates) to read.
	SecretPaths []string
	// RenewToken controls whether the token lease should be renewed automatically.
	RenewToken bool
	// AuditLog is the optional file path for structured audit logging.
	AuditLog string
}

// FromEnv constructs a Config from environment variables.
//
//	VAULT_ADDR        - Vault server address (required)
//	VAULT_TOKEN       - Vault auth token (required)
//	VAULTPIPE_PATHS   - Comma-separated list of secret paths (required)
//	VAULTPIPE_RENEW   - Set to "true" to enable token renewal (optional)
//	VAULTPIPE_AUDIT   - File path for audit log output (optional)
func FromEnv() (*Config, error) {
	cfg := &Config{
		VaultAddr:  os.Getenv("VAULT_ADDR"),
		VaultToken: os.Getenv("VAULT_TOKEN"),
		AuditLog:   os.Getenv("VAULTPIPE_AUDIT"),
		RenewToken: os.Getenv("VAULTPIPE_RENEW") == "true",
	}
	if raw := os.Getenv("VAULTPIPE_PATHS"); raw != "" {
		for _, p := range strings.Split(raw, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				cfg.SecretPaths = append(cfg.SecretPaths, p)
			}
		}
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Validate checks that all required fields are present.
func (c *Config) Validate() error {
	var errs []string
	if c.VaultAddr == "" {
		errs = append(errs, "VAULT_ADDR is required")
	}
	if c.VaultToken == "" {
		errs = append(errs, "VAULT_TOKEN is required")
	}
	if len(c.SecretPaths) == 0 {
		errs = append(errs, "VAULTPIPE_PATHS is required")
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}
