package audit_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/audit"
)

func TestLog_WritesJSON(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	e := audit.Event{
		Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Type:      audit.EventSecretRead,
		Path:      "secret/data/app",
		Success:   true,
	}
	if err := l.Log(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got audit.Event
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("failed to parse output: %v", err)
	}
	if got.Path != "secret/data/app" {
		t.Errorf("expected path secret/data/app, got %q", got.Path)
	}
	if !got.Success {
		t.Error("expected success=true")
	}
}

func TestLog_SetsTimestampIfZero(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	before := time.Now().UTC()
	_ = l.Log(audit.Event{Type: audit.EventExecStart, Success: true})
	after := time.Now().UTC()

	var got audit.Event
	_ = json.Unmarshal(buf.Bytes(), &got)
	if got.Timestamp.Before(before) || got.Timestamp.After(after) {
		t.Errorf("timestamp %v not in expected range", got.Timestamp)
	}
}

func TestSecretRead(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	_ = l.SecretRead("kv/myapp", false)

	var got audit.Event
	_ = json.Unmarshal(buf.Bytes(), &got)
	if got.Type != audit.EventSecretRead {
		t.Errorf("expected type %q, got %q", audit.EventSecretRead, got.Type)
	}
	if got.Success {
		t.Error("expected success=false")
	}
}

func TestExecStartFinish(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	_ = l.ExecStart("env")
	_ = l.ExecFinish("env", true)

	dec := json.NewDecoder(&buf)
	var e1, e2 audit.Event
	_ = dec.Decode(&e1)
	_ = dec.Decode(&e2)
	if e1.Type != audit.EventExecStart {
		t.Errorf("want exec_start, got %q", e1.Type)
	}
	if e2.Type != audit.EventExecFinish {
		t.Errorf("want exec_finish, got %q", e2.Type)
	}
}
