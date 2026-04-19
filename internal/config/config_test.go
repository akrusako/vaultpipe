package config_test

import (
	"testing"

	"github.com/vaultpipe/vaultpipe/internal/config"
)

func TestFromEnv(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://127.0.0.1:8200")
	t.Setenv("VAULT_TOKEN", "test-token")

	cfg := config.FromEnv()
	if cfg.VaultAddr != "http://127.0.0.1:8200" {
		t.Errorf("expected VaultAddr %q, got %q", "http://127.0.0.1:8200", cfg.VaultAddr)
	}
	if cfg.VaultToken != "test-token" {
		t.Errorf("expected VaultToken %q, got %q", "test-token", cfg.VaultToken)
	}
}

func TestValidate_MissingAddr(t *testing.T) {
	cfg := &config.Config{
		VaultToken:  "tok",
		SecretPaths: []string{"secret/foo"},
		Command:     []string{"env"},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for missing VaultAddr")
	}
}

func TestValidate_MissingToken(t *testing.T) {
	cfg := &config.Config{
		VaultAddr:   "http://127.0.0.1:8200",
		SecretPaths: []string{"secret/foo"},
		Command:     []string{"env"},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for missing VaultToken")
	}
}

func TestValidate_MissingSecretPaths(t *testing.T) {
	cfg := &config.Config{
		VaultAddr:  "http://127.0.0.1:8200",
		VaultToken: "tok",
		Command:    []string{"env"},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for missing SecretPaths")
	}
}

func TestValidate_MissingCommand(t *testing.T) {
	cfg := &config.Config{
		VaultAddr:   "http://127.0.0.1:8200",
		VaultToken:  "tok",
		SecretPaths: []string{"secret/foo"},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for missing Command")
	}
}

func TestValidate_Valid(t *testing.T) {
	cfg := &config.Config{
		VaultAddr:   "http://127.0.0.1:8200",
		VaultToken:  "tok",
		SecretPaths: []string{"secret/foo"},
		Command:     []string{"env"},
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
