package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newMockVaultServer(t *testing.T, path string, payload map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": payload,
		})
	}))
}

func TestNewClient_MissingToken(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "")
	_, err := NewClient(Config{Address: "http://127.0.0.1:8200"})
	if err == nil {
		t.Fatal("expected error when token is missing")
	}
}

func TestNewClient_Success(t *testing.T) {
	client, err := NewClient(Config{
		Address: "http://127.0.0.1:8200",
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestFlattenData(t *testing.T) {
	input := map[string]interface{}{
		"DB_PASS": "secret123",
		"PORT":    8080,
	}
	out := flattenData(input)
	if out["DB_PASS"] != "secret123" {
		t.Errorf("expected 'secret123', got %q", out["DB_PASS"])
	}
	if out["PORT"] != "8080" {
		t.Errorf("expected '8080', got %q", out["PORT"])
	}
}

func TestReadSecrets_KVv1Fallback(t *testing.T) {
	server := newMockVaultServer(t, "/v1/secret/myapp", map[string]interface{}{
		"API_KEY": "abc",
	})
	defer server.Close()

	client, err := NewClient(Config{Address: server.URL, Token: "fake"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secrets, err := client.ReadSecrets(context.Background(), "secret/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["API_KEY"] != "abc" {
		t.Errorf("expected 'abc', got %q", secrets["API_KEY"])
	}
}

func TestMountFromPath(t *testing.T) {
	if got := mountFromPath("secret/myapp/config"); got != "secret" {
		t.Errorf("expected 'secret', got %q", got)
	}
}

func TestSubpathFromPath(t *testing.T) {
	if got := subpathFromPath("secret/myapp/config"); got != "myapp/config" {
		t.Errorf("expected 'myapp/config', got %q", got)
	}
}
