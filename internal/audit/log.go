package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Entry represents a single audit log event.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Event     string    `json:"event"`
	Path      string    `json:"path,omitempty"`
	Keys      []string  `json:"keys,omitempty"`
	Error     string    `json:"error,omitempty"`
}

// Logger writes structured audit entries to a writer.
type Logger struct {
	out io.Writer
}

// New creates a Logger writing to the given path.
// Pass an empty path to write to stdout.
func New(logPath string) (*Logger, error) {
	if logPath == "" {
		return &Logger{out: os.Stdout}, nil
	}
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("audit: open log file: %w", err)
	}
	return &Logger{out: f}, nil
}

// NewWithWriter creates a Logger writing to w (useful for testing).
func NewWithWriter(w io.Writer) *Logger {
	return &Logger{out: w}
}

// Log writes an audit entry.
func (l *Logger) Log(e Entry) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}
	_, err = fmt.Fprintln(l.out, string(data))
	return err
}

// LogSync is a convenience method for recording a sync event.
func (l *Logger) LogSync(vaultPath string, keys []string, syncErr error) error {
	e := Entry{
		Event: "sync",
		Path:  vaultPath,
		Keys:  keys,
	}
	if syncErr != nil {
		e.Event = "sync_error"
		e.Error = syncErr.Error()
	}
	return l.Log(e)
}
