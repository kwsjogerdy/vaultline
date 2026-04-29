// Package envtrace records the access history of secret keys,
// tracking when each key was read and by which operation.
package envtrace

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"
)

// Entry represents a single access event for a secret key.
type Entry struct {
	Key       string    `json:"key"`
	Operation string    `json:"operation"`
	At        time.Time `json:"at"`
}

// Tracer records and persists key access traces.
type Tracer struct {
	mu      sync.Mutex
	entries []Entry
	path    string
	clock   func() time.Time
}

// New creates a Tracer that persists entries to path.
// Existing entries are loaded from path if it exists.
func New(path string) (*Tracer, error) {
	t := &Tracer{path: path, clock: time.Now}
	if data, err := os.ReadFile(path); err == nil {
		if err := json.Unmarshal(data, &t.entries); err != nil {
			return nil, fmt.Errorf("envtrace: parse %s: %w", path, err)
		}
	}
	return t, nil
}

// Record appends an access entry for key under the given operation name.
func (t *Tracer) Record(key, operation string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries = append(t.entries, Entry{Key: key, Operation: operation, At: t.clock()})
}

// Entries returns a copy of all recorded entries, sorted by time ascending.
func (t *Tracer) Entries() []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Entry, len(t.entries))
	copy(out, t.entries)
	sort.Slice(out, func(i, j int) bool { return out[i].At.Before(out[j].At) })
	return out
}

// ForKey returns all entries for a specific key.
func (t *Tracer) ForKey(key string) []Entry {
	all := t.Entries()
	var out []Entry
	for _, e := range all {
		if e.Key == key {
			out = append(out, e)
		}
	}
	return out
}

// Save persists all entries to the configured path as JSON.
func (t *Tracer) Save() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	data, err := json.MarshalIndent(t.entries, "", "  ")
	if err != nil {
		return fmt.Errorf("envtrace: marshal: %w", err)
	}
	return os.WriteFile(t.path, data, 0600)
}
