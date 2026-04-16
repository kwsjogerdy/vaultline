package audit

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestLog_WritesJSONEntry(t *testing.T) {
	var buf bytes.Buffer
	l := NewWithWriter(&buf)

	now := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	err := l.Log(Entry{
		Timestamp: now,
		Event:     "sync",
		Path:      "secret/myapp",
		Keys:      []string{"DB_URL", "API_KEY"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var entry Entry
	if err := json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &entry); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if entry.Event != "sync" {
		t.Errorf("expected event=sync, got %q", entry.Event)
	}
	if entry.Path != "secret/myapp" {
		t.Errorf("expected path=secret/myapp, got %q", entry.Path)
	}
	if len(entry.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(entry.Keys))
	}
}

func TestLog_AutoTimestamp(t *testing.T) {
	var buf bytes.Buffer
	l := NewWithWriter(&buf)

	_ = l.Log(Entry{Event: "sync"})

	var entry Entry
	_ = json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &entry)
	if entry.Timestamp.IsZero() {
		t.Error("expected timestamp to be set automatically")
	}
}

func TestLogSync_Success(t *testing.T) {
	var buf bytes.Buffer
	l := NewWithWriter(&buf)

	_ = l.LogSync("secret/app", []string{"FOO", "BAR"}, nil)

	var entry Entry
	_ = json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &entry)
	if entry.Event != "sync" {
		t.Errorf("expected event=sync, got %q", entry.Event)
	}
	if entry.Error != "" {
		t.Errorf("expected no error field, got %q", entry.Error)
	}
}

func TestLogSync_Error(t *testing.T) {
	var buf bytes.Buffer
	l := NewWithWriter(&buf)

	_ = l.LogSync("secret/app", nil, errors.New("unauthorized"))

	var entry Entry
	_ = json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &entry)
	if entry.Event != "sync_error" {
		t.Errorf("expected event=sync_error, got %q", entry.Event)
	}
	if entry.Error != "unauthorized" {
		t.Errorf("expected error=unauthorized, got %q", entry.Error)
	}
}
