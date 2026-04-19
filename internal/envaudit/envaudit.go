// Package envaudit tracks which secrets were accessed, written, or skipped
// during a sync operation, producing a structured access report.
package envaudit

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// Action describes what happened to a secret key.
type Action string

const (
	ActionWritten Action = "written"
	ActionSkipped Action = "skipped"
	ActionRedacted Action = "redacted"
	ActionExpired  Action = "expired"
)

// Entry records a single secret access event.
type Entry struct {
	Key       string    `json:"key"`
	Action    Action    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
	Reason    string    `json:"reason,omitempty"`
}

// Report holds the full audit trail for a sync run.
type Report struct {
	RunID     string    `json:"run_id"`
	StartedAt time.Time `json:"started_at"`
	Entries   []Entry   `json:"entries"`
}

// Tracker accumulates audit entries during a sync.
type Tracker struct {
	report Report
	now    func() time.Time
}

// New creates a new Tracker with the given run ID.
func New(runID string) *Tracker {
	return &Tracker{
		report: Report{
			RunID:     runID,
			StartedAt: time.Now(),
		},
		now: time.Now,
	}
}

// Record adds an audit entry for a key.
func (t *Tracker) Record(key string, action Action, reason string) {
	t.report.Entries = append(t.report.Entries, Entry{
		Key:       key,
		Action:    action,
		Timestamp: t.now(),
		Reason:    reason,
	})
}

// Report returns the completed audit report, with entries sorted by key.
func (t *Tracker) Report() Report {
	sorted := make([]Entry, len(t.report.Entries))
	copy(sorted, t.report.Entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})
	r := t.report
	r.Entries = sorted
	return r
}

// Summary returns a human-readable one-line summary.
func (t *Tracker) Summary() string {
	counts := map[Action]int{}
	for _, e := range t.report.Entries {
		counts[e.Action]++
	}
	return fmt.Sprintf("written=%d skipped=%d redacted=%d expired=%d",
		counts[ActionWritten], counts[ActionSkipped],
		counts[ActionRedacted], counts[ActionExpired])
}

// WriteJSON writes the report as JSON to the given file path.
func (t *Tracker) WriteJSON(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("envaudit: create %s: %w", path, err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(t.Report()); err != nil {
		return fmt.Errorf("envaudit: encode: %w", err)
	}
	return nil
}
