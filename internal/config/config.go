// Package config handles loading and validation of vaultpipe configuration
// from environment variables and CLI flags.
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
	// SecretPaths is the list of Vault secret paths to read.
	SecretPaths []string
	// Command is the process to execute with the injected environment.
	Command []string
}

// FromEnv populates a Config from environment variables.
// VAULT_ADDR and VAULT_TOKEN are read from the environment.
func FromEnv() *Config {
	return &Config{
		VaultAddr:  os.Getenv("VAULT_ADDR"),
		VaultToken: os.Getenv("VAULT_TOKEN"),
	}
}

// Validate checks that the Config has all required fields set.
func (c *Config) Validate() error {
	if strings.TrimSpace(c.VaultAddr) == "" {
		return errors.New("vault address is required (set VAULT_ADDR or --vault-addr)")
	}
	if strings.TrimSpace(c.VaultToken) == "" {
		return errors.New("vault token is required (set VAULT_TOKEN or --vault-token)")
	}
	if len(c.SecretPaths) == 0 {
		return errors.New("at least one secret path must be specified")
	}
	if len(c.Command) == 0 {
		return errors.New("a command to execute is required")
	}
	return nil
}
