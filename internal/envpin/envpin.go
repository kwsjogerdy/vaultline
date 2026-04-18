// Package envpin allows pinning secret versions so that syncs are reproducible.
package envpin

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// Pin records the version metadata for a single secret path.
type Pin struct {
	Path      string    `json:"path"`
	Version   int       `json:"version"`
	PinnedAt  time.Time `json:"pinned_at"`
}

// Lockfile holds all pinned secrets.
type Lockfile struct {
	Pins map[string]Pin `json:"pins"`
}

// Manager manages reading and writing pin lockfiles.
type Manager struct {
	filePath string
}

// New returns a Manager backed by the given file path.
func New(filePath string) *Manager {
	return &Manager{filePath: filePath}
}

// Load reads the lockfile from disk. Returns an empty Lockfile if the file does not exist.
func (m *Manager) Load() (*Lockfile, error) {
	data, err := os.ReadFile(m.filePath)
	if errors.Is(err, os.ErrNotExist) {
		return &Lockfile{Pins: make(map[string]Pin)}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("envpin: read lockfile: %w", err)
	}
	var lf Lockfile
	if err := json.Unmarshal(data, &lf); err != nil {
		return nil, fmt.Errorf("envpin: parse lockfile: %w", err)
	}
	if lf.Pins == nil {
		lf.Pins = make(map[string]Pin)
	}
	return &lf, nil
}

// Save writes the lockfile to disk.
func (m *Manager) Save(lf *Lockfile) error {
	data, err := json.MarshalIndent(lf, "", "  ")
	if err != nil {
		return fmt.Errorf("envpin: marshal lockfile: %w", err)
	}
	if err := os.WriteFile(m.filePath, data, 0600); err != nil {
		return fmt.Errorf("envpin: write lockfile: %w", err)
	}
	return nil
}

// Set pins a path at the given version.
func (m *Manager) Set(lf *Lockfile, path string, version int) {
	lf.Pins[path] = Pin{Path: path, Version: version, PinnedAt: time.Now().UTC()}
}

// Remove deletes a pin entry.
func (m *Manager) Remove(lf *Lockfile, path string) {
	delete(lf.Pins, path)
}

// Get returns the pinned version for a path, or 0 and false if not pinned.
func (m *Manager) Get(lf *Lockfile, path string) (int, bool) {
	p, ok := lf.Pins[path]
	if !ok {
		return 0, false
	}
	return p.Version, true
}
