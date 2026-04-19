package health_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/health"
)

func newHealthServer(code int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
	}))
}

func TestCheck_Active(t *testing.T) {
	srv := newHealthServer(http.StatusOK)
	defer srv.Close()

	c := health.New(srv.URL, 5*time.Second)
	s, err := c.Check(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.IsHealthy() {
		t.Errorf("expected healthy, got %+v", s)
	}
}

func TestCheck_Standby(t *testing.T) {
	srv := newHealthServer(429)
	defer srv.Close()

	c := health.New(srv.URL, 5*time.Second)
	s, err := c.Check(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.IsHealthy() {
		t.Error("standby node should not be healthy")
	}
	if !s.Standby {
		t.Error("expected Standby=true")
	}
}

func TestCheck_Sealed(t *testing.T) {
	srv := newHealthServer(503)
	defer srv.Close()

	c := health.New(srv.URL, 5*time.Second)
	s, err := c.Check(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.Sealed {
		t.Error("expected Sealed=true")
	}
}

func TestCheck_UnexpectedCode(t *testing.T) {
	srv := newHealthServer(http.StatusTeapot)
	defer srv.Close()

	c := health.New(srv.URL, 5*time.Second)
	_, err := c.Check(context.Background())
	if err == nil {
		t.Error("expected error for unexpected status code")
	}
}

func TestCheck_Unreachable(t *testing.T) {
	c := health.New("http://127.0.0.1:1", 1*time.Second)
	_, err := c.Check(context.Background())
	if err == nil {
		t.Error("expected error for unreachable server")
	}
}
