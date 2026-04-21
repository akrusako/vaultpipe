// Package notify provides a lightweight event notification system for
// vaultpipe. It allows internal components to emit informational, warning,
// and error events without being coupled to a specific logging backend.
//
// Consumers register Handler functions that receive Event values. When no
// handlers are registered the notifier falls back to writing formatted lines
// to a configured io.Writer (defaulting to os.Stderr).
//
// Usage:
//
//	n := notify.New(os.Stderr)
//	n.Register(func(ctx context.Context, e notify.Event) {
//		log.Printf("[%s] %s", e.Level, e.Message)
//	})
//	n.Warn(ctx, "secret lease expiring soon")
package notify
