package rotate_test

import (
	"context"
	"errors"
	"log/slog"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/rotate"
)

type mockFetcher struct {
	calls   atomic.Int32
	secrets map[string]string
	err     error
}

func (m *mockFetcher) ReadSecrets(_ context.Context, _ []string) (map[string]string, error) {
	m.calls.Add(1)
	return m.secrets, m.err
}

func TestRotator_CallsHandlerOnTick(t *testing.T) {
	fetcher := &mockFetcher{secrets: map[string]string{"KEY": "val"}}
	var received atomic.Int32

	r := rotate.New(fetcher, []string{"secret/app"}, 20*time.Millisecond,
		func(_ map[string]string) { received.Add(1) },
		slog.Default())

	ctx, cancel := context.WithTimeout(context.Background(), 70*time.Millisecond)
	defer cancel()
	r.Run(ctx)

	if received.Load() < 2 {
		t.Fatalf("expected at least 2 handler calls, got %d", received.Load())
	}
}

func TestRotator_SkipsHandlerOnError(t *testing.T) {
	fetcher := &mockFetcher{err: errors.New("vault down")}
	var received atomic.Int32

	r := rotate.New(fetcher, []string{"secret/app"}, 20*time.Millisecond,
		func(_ map[string]string) { received.Add(1) },
		nil)

	ctx, cancel := context.WithTimeout(context.Background(), 70*time.Millisecond)
	defer cancel()
	r.Run(ctx)

	if received.Load() != 0 {
		t.Fatalf("expected 0 handler calls on error, got %d", received.Load())
	}
	if fetcher.calls.Load() < 2 {
		t.Fatalf("expected fetcher to be called despite errors")
	}
}

func TestRotator_StopsOnContextCancel(t *testing.T) {
	fetcher := &mockFetcher{secrets: map[string]string{}}
	ctx, cancel := context.WithCancel(context.Background())

	r := rotate.New(fetcher, nil, 10*time.Millisecond,
		func(_ map[string]string) {},
		slog.Default())

	done := make(chan struct{})
	go func() { r.Run(ctx); close(done) }()

	cancel()
	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("rotator did not stop after context cancel")
	}
}
