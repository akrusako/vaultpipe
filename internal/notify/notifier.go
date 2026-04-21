package notify

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the severity of a notification.
type Level int

const (
	LevelInfo Level = iota
	LevelWarn
	LevelError
)

func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Event holds the data for a single notification.
type Event struct {
	Level   Level
	Message string
	At      time.Time
}

// Handler is a function that receives a notification event.
type Handler func(ctx context.Context, e Event)

// Notifier dispatches events to registered handlers.
type Notifier struct {
	handlers []Handler
	out      io.Writer
}

// New creates a Notifier that writes to out when no handlers are registered.
func New(out io.Writer) *Notifier {
	if out == nil {
		out = os.Stderr
	}
	return &Notifier{out: out}
}

// Register adds a handler to the notifier.
func (n *Notifier) Register(h Handler) {
	n.handlers = append(n.handlers, h)
}

// Send dispatches an event at the given level to all registered handlers.
// If no handlers are registered it falls back to writing to the configured writer.
func (n *Notifier) Send(ctx context.Context, level Level, msg string) {
	e := Event{Level: level, Message: msg, At: time.Now().UTC()}
	if len(n.handlers) == 0 {
		fmt.Fprintf(n.out, "[%s] %s %s\n", e.Level, e.At.Format(time.RFC3339), e.Message)
		return
	}
	for _, h := range n.handlers {
		h(ctx, e)
	}
}

// Info is a convenience wrapper for LevelInfo events.
func (n *Notifier) Info(ctx context.Context, msg string) {
	n.Send(ctx, LevelInfo, msg)
}

// Warn is a convenience wrapper for LevelWarn events.
func (n *Notifier) Warn(ctx context.Context, msg string) {
	n.Send(ctx, LevelWarn, msg)
}

// Error is a convenience wrapper for LevelError events.
func (n *Notifier) Error(ctx context.Context, msg string) {
	n.Send(ctx, LevelError, msg)
}
