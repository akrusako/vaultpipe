// Package audit provides structured audit logging for secret access events.
package audit

import (
	"encoding/json"
	"io"
	"time"
)

// EventType classifies the kind of audit event.
type EventType string

const (
	EventSecretRead  EventType = "secret_read"
	EventExecStart   EventType = "exec_start"
	EventExecFinish  EventType = "exec_finish"
	EventRenewToken  EventType = "renew_token"
	EventRenewFailed EventType = "renew_failed"
)

// Event represents a single auditable action.
type Event struct {
	Timestamp time.Time  `json:"timestamp"`
	Type      EventType  `json:"type"`
	Path      string     `json:"path,omitempty"`
	Command   string     `json:"command,omitempty"`
	Message   string     `json:"message,omitempty"`
	Success   bool       `json:"success"`
}

// Logger writes audit events as newline-delimited JSON.
type Logger struct {
	w   io.Writer
	enc *json.Encoder
}

// New creates a new audit Logger writing to w.
func New(w io.Writer) *Logger {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return &Logger{w: w, enc: enc}
}

// Log emits an audit event.
func (l *Logger) Log(e Event) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	return l.enc.Encode(e)
}

// SecretRead logs a secret read event for the given path.
func (l *Logger) SecretRead(path string, ok bool) error {
	return l.Log(Event{Type: EventSecretRead, Path: path, Success: ok})
}

// ExecStart logs that a child process is about to start.
func (l *Logger) ExecStart(cmd string) error {
	return l.Log(Event{Type: EventExecStart, Command: cmd, Success: true})
}

// ExecFinish logs that a child process finished.
func (l *Logger) ExecFinish(cmd string, ok bool) error {
	return l.Log(Event{Type: EventExecFinish, Command: cmd, Success: ok})
}
