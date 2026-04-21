package notify_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/yourusername/vaultpipe/internal/notify"
)

func TestNew_DefaultsToStderr(t *testing.T) {
	n := notify.New(nil)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestSend_FallbackWriter(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(&buf)
	n.Send(context.Background(), notify.LevelInfo, "hello world")

	if !strings.Contains(buf.String(), "hello world") {
		t.Errorf("expected output to contain message, got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "INFO") {
		t.Errorf("expected output to contain level, got: %s", buf.String())
	}
}

func TestSend_CallsHandler(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(&buf)

	var received notify.Event
	n.Register(func(_ context.Context, e notify.Event) {
		received = e
	})

	n.Send(context.Background(), notify.LevelWarn, "rotation due")

	if received.Message != "rotation due" {
		t.Errorf("expected message 'rotation due', got %q", received.Message)
	}
	if received.Level != notify.LevelWarn {
		t.Errorf("expected LevelWarn, got %v", received.Level)
	}
	if buf.Len() != 0 {
		t.Error("expected fallback writer not to be used when handler registered")
	}
}

func TestSend_MultipleHandlers(t *testing.T) {
	n := notify.New(nil)
	calls := 0
	for i := 0; i < 3; i++ {
		n.Register(func(_ context.Context, _ notify.Event) { calls++ })
	}
	n.Info(context.Background(), "multi")
	if calls != 3 {
		t.Errorf("expected 3 handler calls, got %d", calls)
	}
}

func TestLevelString(t *testing.T) {
	cases := []struct {
		level notify.Level
		want  string
	}{
		{notify.LevelInfo, "INFO"},
		{notify.LevelWarn, "WARN"},
		{notify.LevelError, "ERROR"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("Level.String() = %q, want %q", got, tc.want)
		}
	}
}

func TestConvenienceMethods(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(&buf)

	levels := []notify.Level{}
	n.Register(func(_ context.Context, e notify.Event) {
		levels = append(levels, e.Level)
	})

	n.Info(context.Background(), "a")
	n.Warn(context.Background(), "b")
	n.Error(context.Background(), "c")

	if len(levels) != 3 {
		t.Fatalf("expected 3 events, got %d", len(levels))
	}
	if levels[0] != notify.LevelInfo || levels[1] != notify.LevelWarn || levels[2] != notify.LevelError {
		t.Errorf("unexpected level sequence: %v", levels)
	}
}
