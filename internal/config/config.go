// Package config loads vaultpipe runtime configuration from environment variables.
package config

import (
	"errors"
	"os"
	"strings"
)

// Config holds all runtime configuration for vaultpipe.
type Config struct {
	// VaultAddr is the Vault server address.
	VaultAddr string
	// VaultToken is the Vault authentication token.
	VaultToken string
	// SecretPaths is the list of KV paths to read.
	SecretPaths []string
	// Command is the process to execute with secrets in its environment.
	Command []string
	// AuditLog enables structured audit logging to stderr when true.
	AuditLog bool
	// LogLevel controls verbosity (debug, info, warn, error).
	LogLevel string
}

// FromEnv reads configuration from environment variables.
func FromEnv() (*Config, error) {
	cfg := &Config{
		VaultAddr:  os.Getenv("VAULT_ADDR"),
		VaultToken: os.Getenv("VAULT_TOKEN"),
		LogLevel:   os.Getenv("VAULTPIPE_LOG_LEVEL"),
		AuditLog:   os.Getenv("VAULTPIPE_AUDIT") == "true",
	}

	if raw := os.Getenv("VAULTPIPE_SECRET_PATHS"); raw != "" {
		for _, p := range strings.Split(raw, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				cfg.SecretPaths = append(cfg.SecretPaths, p)
			}
		}
	}

	if raw := os.Getenv("VAULTPIPE_COMMAND"); raw != "" {
		cfg.Command = strings.Fields(raw)
	}

	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}

	return cfg, nil
}

// Validate returns an error if any required fields are missing.
func (c *Config) Validate() error {
	if c.VaultAddr == "" {
		return errors.New("VAULT_ADDR is required")
	}
	if c.VaultToken == "" {
		return errors.New("VAULT_TOKEN is required")
	}
	if len(c.SecretPaths) == 0 {
		return errors.New("VAULTPIPE_SECRET_PATHS is required")
	}
	if len(c.Command) == 0 {
		return errors.New("VAULTPIPE_COMMAND is required")
	}
	return nil
}
