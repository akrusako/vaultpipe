// Package log provides a lightweight structured logger for vaultpipe.
// It wraps the standard library's log/slog package and adds support for
// secret masking so sensitive values are never emitted to output.
package log

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/yourusername/vaultpipe/internal/mask"
)

// Logger wraps slog.Logger with an integrated secret masker.
type Logger struct {
	inner  *slog.Logger
	masker *mask.Masker
}

// Level represents the logging verbosity.
type Level = slog.Level

const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

// New creates a Logger that writes JSON-structured output to w at the given
// level. If w is nil, os.Stderr is used.
func New(w io.Writer, level Level) *Logger {
	if w == nil {
		w = os.Stderr
	}
	h := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: level,
	})
	return &Logger{
		inner:  slog.New(h),
		masker: mask.New(nil),
	}
}

// SetSecrets registers secret values that must be redacted from all log output.
func (l *Logger) SetSecrets(secrets map[string]string) {
	for _, v := range secrets {
		l.masker.Add(v)
	}
}

// redact applies masking to every string argument in args.
func (l *Logger) redact(msg string, args ...any) (string, []any) {
	msg = l.masker.Redact(msg)
	for i, a := range args {
		if s, ok := a.(string); ok {
			args[i] = l.masker.Redact(s)
		}
	}
	return msg, args
}

// Debug logs at DEBUG level.
func (l *Logger) Debug(msg string, args ...any) {
	msg, args = l.redact(msg, args...)
	l.inner.Debug(msg, args...)
}

// Info logs at INFO level.
func (l *Logger) Info(msg string, args ...any) {
	msg, args = l.redact(msg, args...)
	l.inner.Info(msg, args...)
}

// Warn logs at WARN level.
func (l *Logger) Warn(msg string, args ...any) {
	msg, args = l.redact(msg, args...)
	l.inner.Warn(msg, args...)
}

// Error logs at ERROR level.
func (l *Logger) Error(msg string, args ...any) {
	msg, args = l.redact(msg, args...)
	l.inner.Error(msg, args...)
}

// With returns a new Logger with the given attributes pre-attached.
func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		inner:  l.inner.With(args...),
		masker: l.masker,
	}
}

// InfoContext logs at INFO level with a context.
func (l *Logger) InfoContext(ctx context.Context, msg string, args ...any) {
	msg, args = l.redact(msg, args...)
	l.inner.InfoContext(ctx, msg, args...)
}

// ErrorContext logs at ERROR level with a context.
func (l *Logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	msg, args = l.redact(msg, args...)
	l.inner.ErrorContext(ctx, msg, args...)
}
