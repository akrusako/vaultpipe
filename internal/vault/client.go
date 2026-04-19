package vault

import (
	"context"
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with helper methods.
type Client struct {
	vc *vaultapi.Client
}

// Config holds configuration for connecting to Vault.
type Config struct {
	Address string
	Token   string
}

// NewClient creates a new Vault client from the given config.
// If Address or Token are empty, it falls back to environment variables
// VAULT_ADDR and VAULT_TOKEN respectively.
func NewClient(cfg Config) (*Client, error) {
	vcfg := vaultapi.DefaultConfig()

	addr := cfg.Address
	if addr == "" {
		addr = os.Getenv("VAULT_ADDR")
	}
	if addr == "" {
		addr = "http://127.0.0.1:8200"
	}
	vcfg.Address = addr

	vc, err := vaultapi.NewClient(vcfg)
	if err != nil {
		return nil, fmt.Errorf("vault: failed to create client: %w", err)
	}

	token := cfg.Token
	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}
	if token == "" {
		return nil, fmt.Errorf("vault: no token provided (set VAULT_TOKEN or pass --token)")
	}
	vc.SetToken(token)

	return &Client{vc: vc}, nil
}

// ReadSecrets reads the secret at the given KV path and returns a map of
// key/value pairs. Supports both KV v1 and v2 mounts.
func (c *Client) ReadSecrets(ctx context.Context, path string) (map[string]string, error) {
	secret, err := c.vc.KVv2(mountFromPath(path)).Get(ctx, subpathFromPath(path))
	if err != nil {
		// Fallback to KV v1
		raw, err2 := c.vc.Logical().ReadWithContext(ctx, path)
		if err2 != nil {
			return nil, fmt.Errorf("vault: read %q failed: %w", path, err)
		}
		if raw == nil || raw.Data == nil {
			return nil, fmt.Errorf("vault: no data at path %q", path)
		}
		return flattenData(raw.Data), nil
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("vault: no data at path %q", path)
	}
	return flattenData(secret.Data), nil
}

func mountFromPath(path string) string {
	for i, c := range path {
		if c == '/' {
			return path[:i]
		}
	}
	return path
}

func subpathFromPath(path string) string {
	for i, c := range path {
		if c == '/' {
			return path[i+1:]
		}
	}
	return ""
}

func flattenData(data map[string]interface{}) map[string]string {
	out := make(map[string]string, len(data))
	for k, v := range data {
		out[k] = fmt.Sprintf("%v", v)
	}
	return out
}
