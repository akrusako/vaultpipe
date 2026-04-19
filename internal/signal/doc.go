// Package signal provides OS signal handling for vaultpipe.
//
// It wraps the standard library's signal.Notify to allow graceful shutdown
// by cancelling a context when SIGINT or SIGTERM is received. This enables
// child processes and background goroutines (such as token renewers) to
// clean up before the process exits.
//
// Usage:
//
//	h := signal.New()
//	ctx, cancel := h.Watch(context.Background())
//	defer cancel()
package signal
