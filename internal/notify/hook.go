package notify

import (
	"context"
	"fmt"
	"io"
	"time"
)

// JSONHook returns a Handler that writes newline-delimited JSON events to w.
// It is suitable for structured logging pipelines.
func JSONHook(w io.Writer) Handler {
	return func(_ context.Context, e Event) {
		fmt.Fprintf(w, `{"level":%q,"message":%q,"at":%q}`+"\n",
			e.Level.String(),
			e.Message,
			e.At.Format(time.RFC3339),
		)
	}
}

// FilterHook wraps an existing Handler and only forwards events whose level
// is greater than or equal to min.
func FilterHook(min Level, next Handler) Handler {
	return func(ctx context.Context, e Event) {
		if e.Level >= min {
			next(ctx, e)
		}
	}
}

// ChannelHook returns a Handler that sends events to ch in a non-blocking
// fashion. Events are dropped if ch is full.
func ChannelHook(ch chan<- Event) Handler {
	return func(_ context.Context, e Event) {
		select {
		case ch <- e:
		default:
		}
	}
}
