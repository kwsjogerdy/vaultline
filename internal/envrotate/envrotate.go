// Package envrotate provides secret rotation tracking for .env files.
package envrotate

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Record holds rotation metadata for a single secret key.
type Record struct {
	Key       string    `json:"key"`
	RotatedAt time.Time `json:"rotated_at"`
	Version   int       `json:"version"`
}

// Ledger tracks rotation history for secrets.
type Ledger struct {
	path    string
	records map[string]Record
}

// New loads or initialises a rotation ledger at path.
func New(path string) (*Ledger, error) {
	l := &Ledger{path: path, records: make(map[string]Record)}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return l, nil
	}
	if err != nil {
		return nil, fmt.Errorf("envrotate: read ledger: %w", err)
	}
	if err := json.Unmarshal(data, &l.records); err != nil {
		return nil, fmt.Errorf("envrotate: parse ledger: %w", err)
	}
	return l, nil
}

// Rotate marks key as rotated now, incrementing its version.
func (l *Ledger) Rotate(key string) Record {
	r := l.records[key]
	r.Key = key
	r.RotatedAt = time.Now().UTC()
	r.Version++
	l.records[key] = r
	return r
}

// Get returns the rotation record for key, and whether it exists.
func (l *Ledger) Get(key string) (Record, bool) {
	r, ok := l.records[key]
	return r, ok
}

// DueForRotation returns keys whose last rotation is older than maxAge.
func (l *Ledger) DueForRotation(keys []string, maxAge time.Duration) []string {
	var due []string
	now := time.Now().UTC()
	for _, k := range keys {
		r, ok := l.records[k]
		if !ok || now.Sub(r.RotatedAt) > maxAge {
			due = append(due, k)
		}
	}
	return due
}

// Save persists the ledger to disk.
func (l *Ledger) Save() error {
	data, err := json.MarshalIndent(l.records, "", "  ")
	if err != nil {
		return fmt.Errorf("envrotate: marshal: %w", err)
	}
	if err := os.WriteFile(l.path, data, 0600); err != nil {
		return fmt.Errorf("envrotate: write ledger: %w", err)
	}
	return nil
}
