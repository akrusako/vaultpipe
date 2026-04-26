// Package observe provides a composable pipeline observer that aggregates
// metrics, audit events, and notifications into a single coordinated hook
// invoked at key points in the secret-fetch lifecycle.
package observe

import (
	"context"
	"time"

	"github.com/yourusername/vaultpipe/internal/audit"
	"github.com/yourusername/vaultpipe/internal/metrics"
	"github.com/yourusername/vaultpipe/internal/notify"
)

// Event describes a lifecycle moment in the vaultpipe pipeline.
type Event struct {
	// Path is the Vault secret path involved, if any.
	Path string

	// Keys is the list of secret keys that were read or changed.
	Keys []string

	// Err holds any error associated with the event, or nil on success.
	Err error

	// Duration is the elapsed time for the operation, if measured.
	Duration time.Duration

	// At is the wall-clock time of the event; set automatically if zero.
	At time.Time
}

// Observer coordinates audit logging, metrics recording, and user
// notifications for pipeline lifecycle events.
type Observer struct {
	auditor  *audit.Logger
	metrics  *metrics.Registry
	notifier *notify.Notifier
}

// New constructs an Observer.  All three dependencies are required; pass
// no-op instances if a subsystem is not needed.
func New(a *audit.Logger, m *metrics.Registry, n *notify.Notifier) *Observer {
	return &Observer{
		auditor:  a,
		metrics:  m,
		notifier: n,
	}
}

// SecretFetched records a successful secret read.  It increments the
// "secrets_fetched" counter, writes an audit entry, and emits an info
// notification with the path and key count.
func (o *Observer) SecretFetched(ctx context.Context, ev Event) {
	ev = normalise(ev)

	o.metrics.Inc("secrets_fetched")
	if ev.Duration > 0 {
		o.metrics.Inc("fetch_duration_ms")
	}

	o.auditor.SecretRead(ev.Path, ev.Keys)

	o.notifier.Send(ctx, notify.Info,
		"secret fetched path=%s keys=%d dur=%s",
		ev.Path, len(ev.Keys), ev.Duration.Round(time.Millisecond))
}

// SecretChanged records that a watched secret changed on a poll cycle.
// It increments "secrets_changed" and emits a warning-level notification
// so operators are aware of rotation events.
func (o *Observer) SecretChanged(ctx context.Context, ev Event) {
	ev = normalise(ev)

	o.metrics.Inc("secrets_changed")
	o.auditor.SecretRead(ev.Path, ev.Keys)

	o.notifier.Send(ctx, notify.Warn,
		"secret changed path=%s keys=%d",
		ev.Path, len(ev.Keys))
}

// FetchError records a failed secret read.  It increments "fetch_errors"
// and emits an error-level notification containing the path and error text.
func (o *Observer) FetchError(ctx context.Context, ev Event) {
	ev = normalise(ev)

	o.metrics.Inc("fetch_errors")

	msg := "fetch error path=%s"
	if ev.Err != nil {
		msg += " err=" + ev.Err.Error()
	}
	o.notifier.Send(ctx, notify.Error, msg, ev.Path)
}

// ExecStarted records the moment a child process is launched.  It
// increments "exec_started" and writes an audit exec-start entry.
func (o *Observer) ExecStarted(ctx context.Context, cmd string) {
	o.metrics.Inc("exec_started")
	o.auditor.ExecStart(cmd)
}

// ExecFinished records the completion of a child process.  It increments
// "exec_finished" (or "exec_errors" on non-nil err) and writes an audit
// exec-finish entry.
func (o *Observer) ExecFinished(ctx context.Context, cmd string, err error) {
	if err != nil {
		o.metrics.Inc("exec_errors")
	} else {
		o.metrics.Inc("exec_finished")
	}
	o.auditor.ExecFinish(cmd, err)
}

// normalise fills zero-value fields so callers don't have to.
func normalise(ev Event) Event {
	if ev.At.IsZero() {
		ev.At = time.Now().UTC()
	}
	return ev
}
