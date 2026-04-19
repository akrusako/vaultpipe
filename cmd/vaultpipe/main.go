package main

import (
	"fmt"
	"log"
	"os"

	"github.com/yourusername/vaultpipe/internal/config"
	"github.com/yourusername/vaultpipe/internal/exec"
	"github.com/yourusername/vaultpipe/internal/vault"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("vaultpipe %s\n", version)
		os.Exit(0)
	}

	cfg, err := config.FromEnv()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid config: %v", err)
	}

	client, err := vault.NewClient(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		log.Fatalf("vault client error: %v", err)
	}

	secrets := make(map[string]string)
	for _, path := range cfg.SecretPaths {
		data, err := client.ReadSecrets(path)
		if err != nil {
			log.Fatalf("failed to read secrets from %s: %v", path, err)
		}
		for k, v := range data {
			secrets[k] = v
		}
	}

	runner := exec.NewRunner(cfg.Command, secrets)
	if err := runner.Run(); err != nil {
		log.Fatalf("process error: %v", err)
	}
}
