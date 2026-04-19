// Package signal provides graceful shutdown handling for vaultpipe,
// forwarding OS signals to child processes and triggering cleanup.
package signal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// Handler listens for OS signals and cancels a context.
type Handler struct {
	signals []os.Signal
	notify  func(chan<- os.Signal, ...os.Signal)
	stop    func(chan<- os.Signal)
}

// New returns a Handler that watches SIGINT and SIGTERM by default.
func New(extra ...os.Signal) *Handler {
	sigs := []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	sigs = append(sigs, extra...)
	return &Handler{
		signals: sigs,
		notify:  signal.Notify,
		stop:    signal.Stop,
	}
}

// Watch starts a goroutine that cancels ctx when a signal is received.
// The returned context is derived from the provided parent.
func (h *Handler) Watch(parent context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)
	ch := make(chan os.Signal, 1)
	h.notify(ch, h.signals...)

	go func() {
		defer h.stop(ch)
		select {
		case <-ch:
			cancel()
		case <-parent.Done():
			cancel()
		}
	}()

	return ctx, cancel
}
