package metrics

import (
	"context"
	"encoding/json"
	"io"
	"time"
)

// Emitter periodically writes a JSON snapshot of a Metrics registry to a
// writer (e.g. a structured log stream or debug socket).
type Emitter struct {
	m        *Metrics
	w        io.Writer
	interval time.Duration
}

// NewEmitter returns an Emitter that will write to w every interval.
// A zero or negative interval defaults to 60 seconds.
func NewEmitter(m *Metrics, w io.Writer, interval time.Duration) *Emitter {
	if interval <= 0 {
		interval = 60 * time.Second
	}
	return &Emitter{m: m, w: w, interval: interval}
}

// Run blocks, emitting snapshots until ctx is cancelled.
func (e *Emitter) Run(ctx context.Context) {
	ticker := time.NewTicker(e.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			e.emit(t)
		}
	}
}

type emitPayload struct {
	Timestamp string            `json:"timestamp"`
	Counters  map[string]uint64 `json:"counters"`
}

func (e *Emitter) emit(t time.Time) {
	payload := emitPayload{
		Timestamp: t.UTC().Format(time.RFC3339),
		Counters:  e.m.Snapshot(),
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}
	data = append(data, '\n')
	_, _ = e.w.Write(data)
}
