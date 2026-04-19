package config_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/config"
)

func setEnv(t *testing.T, kv map[string]string) {
	t.Helper()
	for k, v := range kv {
		t.Setenv(k, v)
	}
}

func TestFromEnv(t *testing.T) {
	setEnv(t, map[string]string{
		"VAULT_ADDR":             "http://127.0.0.1:8200",
		"VAULT_TOKEN":            "root",
		"VAULTPIPE_SECRET_PATHS": "secret/app, secret/db",
		"VAULTPIPE_COMMAND":      "env -i",
		"VAULTPIPE_AUDIT":        "true",
		"VAULTPIPE_LOG_LEVEL":    "debug",
	})
	cfg, err := config.FromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.VaultAddr != "http://127.0.0.1:8200" {
		t.Errorf("unexpected addr: %q", cfg.VaultAddr)
	}
	if len(cfg.SecretPaths) != 2 {
		t.Errorf("expected 2 paths, got %d", len(cfg.SecretPaths))
	}
	if !cfg.AuditLog {
		t.Error("expected AuditLog=true")
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("expected log level debug, got %q", cfg.LogLevel)
	}
}

func TestValidate_MissingAddr(t *testing.T) {
	cfg := &config.Config{VaultToken: "x", SecretPaths: []string{"s"}, Command: []string{"cmd"}}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing addr")
	}
}

func TestValidate_MissingToken(t *testing.T) {
	cfg := &config.Config{VaultAddr: "http://x", SecretPaths: []string{"s"}, Command: []string{"cmd"}}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing token")
	}
}

func TestValidate_MissingSecretPaths(t *testing.T) {
	cfg := &config.Config{VaultAddr: "http://x", VaultToken: "t", Command: []string{"cmd"}}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing paths")
	}
}

func TestValidate_MissingCommand(t *testing.T) {
	cfg := &config.Config{VaultAddr: "http://x", VaultToken: "t", SecretPaths: []string{"s"}}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing command")
	}
}

func TestValidate_Valid(t *testing.T) {
	cfg := &config.Config{
		VaultAddr:   "http://x",
		VaultToken:  "t",
		SecretPaths: []string{"s"},
		Command:     []string{"cmd"},
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
