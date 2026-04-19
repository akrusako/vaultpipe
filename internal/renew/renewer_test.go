package renew_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	vaultapi "github.com/hashicorp/vault/api"

	"github.com/yourusername/vaultpipe/internal/renew"
)

type stubLogger struct {
	infos []string
	warns []string
}

func (s *stubLogger) Infof(f string, a ...interface{}) { s.infos = append(s.infos, f) }
func (s *stubLogger) Warnf(f string, a ...interface{}) { s.warns = append(s.warns, f) }

func newTestVaultClient(t *testing.T, handler http.Handler) *vaultapi.Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	cfg := vaultapi.DefaultConfig()
	cfg.Address = srv.URL
	c, err := vaultapi.NewClient(cfg)
	if err != nil {
		t.Fatalf("vault client: %v", err)
	}
	c.SetToken("test-token")
	return c
}

func TestRenewer_RenewsCalled(t *testing.T) {
	var calls atomic.Int32
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/auth/token/renew-self" {
			calls.Add(1)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"auth":{"client_token":"test-token","lease_duration":3600}}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	})

	client := newTestVaultClient(t, h)
	log := &stubLogger{}
	r := renew.New(client, 50*time.Millisecond, log)

	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Millisecond)
	defer cancel()
	r.Start(ctx)

	if calls.Load() < 2 {
		t.Errorf("expected at least 2 renewal calls, got %d", calls.Load())
	}
}

func TestRenewer_WarnOnFailure(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})
	client := newTestVaultClient(t, h)
	log := &stubLogger{}
	r := renew.New(client, 30*time.Millisecond, log)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	r.Start(ctx)

	if len(log.warns) == 0 {
		t.Error("expected warnings on renewal failure")
	}
}
