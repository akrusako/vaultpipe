package flatten_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/flatten"
)

func TestFlatten_EmptyInput(t *testing.T) {
	f := flatten.New(flatten.DefaultOptions())
	out := f.Flatten(map[string]any{})
	if len(out) != 0 {
		t.Fatalf("expected empty map, got %v", out)
	}
}

func TestFlatten_SimpleKeys(t *testing.T) {
	f := flatten.New(flatten.DefaultOptions())
	out := f.Flatten(map[string]any{
		"username": "admin",
		"password": "s3cr3t",
	})
	if out["USERNAME"] != "admin" {
		t.Errorf("expected USERNAME=admin, got %q", out["USERNAME"])
	}
	if out["PASSWORD"] != "s3cr3t" {
		t.Errorf("expected PASSWORD=s3cr3t, got %q", out["PASSWORD"])
	}
}

func TestFlatten_NestedMap(t *testing.T) {
	f := flatten.New(flatten.DefaultOptions())
	out := f.Flatten(map[string]any{
		"db": map[string]any{
			"host": "localhost",
			"port": "5432",
		},
	})
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", out["DB_HOST"])
	}
	if out["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", out["DB_PORT"])
	}
}

func TestFlatten_StringStringMap(t *testing.T) {
	f := flatten.New(flatten.DefaultOptions())
	out := f.Flatten(map[string]any{
		"creds": map[string]string{
			"api_key": "abc123",
		},
	})
	if out["CREDS_API_KEY"] != "abc123" {
		t.Errorf("expected CREDS_API_KEY=abc123, got %q", out["CREDS_API_KEY"])
	}
}

func TestFlatten_CustomSeparator(t *testing.T) {
	f := flatten.New(flatten.Options{Separator: ".", Uppercase: false})
	out := f.Flatten(map[string]any{
		"db": map[string]any{
			"host": "127.0.0.1",
		},
	})
	if out["db.host"] != "127.0.0.1" {
		t.Errorf("expected db.host=127.0.0.1, got %q", out["db.host"])
	}
}

func TestFlatten_NonStringValue(t *testing.T) {
	f := flatten.New(flatten.DefaultOptions())
	out := f.Flatten(map[string]any{
		"retries": 3,
	})
	if out["RETRIES"] != "3" {
		t.Errorf("expected RETRIES=3, got %q", out["RETRIES"])
	}
}

func TestDefaultOptions_Separator(t *testing.T) {
	opts := flatten.DefaultOptions()
	if opts.Separator != "_" {
		t.Errorf("expected separator '_', got %q", opts.Separator)
	}
	if !opts.Uppercase {
		t.Error("expected Uppercase to be true")
	}
}
