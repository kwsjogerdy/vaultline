package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp string `json:"timestamp,omitempty"`
	Event     string `json:"event"`
	Detail    string `json:"detail,omitempty"`
	Error     string `json:"error,omitempty"`
}

// Logger writes structured audit entries.
type Logger struct {
	w io.Writer
}

// New creates a Logger that appends to the given file path.
func New(path string) (*Logger, error) {
	if path == "" {
		return &Logger{w: io.Discard}, nil
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("open audit log: %w", err)
	}
	return &Logger{w: f}, nil
}

// NewWithWriter creates a Logger writing to the provided writer (useful for tests).
func NewWithWriter(w io.Writer) *Logger {
	return &Logger{w: w}
}

// Log writes a single Entry, stamping the timestamp if absent.
func (l *Logger) Log(e Entry) error {
	if e.Timestamp == "" {
		e.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(l.w, "%s\n", data)
	return err
}

// LogSync is a convenience wrapper for logging a sync operation result.
func (l *Logger) LogSync(path string, count int, err error) error {
	e := Entry{Event: "sync", Detail: fmt.Sprintf("%s (%d keys)", path, count)}
	if err != nil {
		e.Error = err.Error()
	}
	return l.Log(e)
}
