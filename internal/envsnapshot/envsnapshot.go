// Package envsnapshot provides point-in-time snapshots of secret maps.
package envsnapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Snapshot represents a saved state of secrets at a point in time.
type Snapshot struct {
	CreatedAt time.Time         `json:"created_at"`
	Label     string            `json:"label"`
	Secrets   map[string]string `json:"secrets"`
}

// Manager handles saving and loading snapshots.
type Manager struct {
	dir string
}

// New returns a Manager storing snapshots in dir.
func New(dir string) *Manager {
	return &Manager{dir: dir}
}

// Save writes secrets as a named snapshot to disk.
func (m *Manager) Save(label string, secrets map[string]string) error {
	if err := os.MkdirAll(m.dir, 0700); err != nil {
		return fmt.Errorf("envsnapshot: mkdir: %w", err)
	}
	snap := Snapshot{
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Secrets:   secrets,
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("envsnapshot: marshal: %w", err)
	}
	name := fmt.Sprintf("%d_%s.json", snap.CreatedAt.UnixNano(), label)
	return os.WriteFile(filepath.Join(m.dir, name), data, 0600)
}

// List returns all snapshots sorted newest first.
func (m *Manager) List() ([]Snapshot, error) {
	entries, err := os.ReadDir(m.dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("envsnapshot: readdir: %w", err)
	}
	var snaps []Snapshot
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(m.dir, e.Name()))
		if err != nil {
			continue
		}
		var s Snapshot
		if json.Unmarshal(data, &s) == nil {
			snaps = append(snaps, s)
		}
	}
	sort.Slice(snaps, func(i, j int) bool {
		return snaps[i].CreatedAt.After(snaps[j].CreatedAt)
	})
	return snaps, nil
}

// Latest returns the most recently saved snapshot, or nil if none exist.
func (m *Manager) Latest() (*Snapshot, error) {
	list, err := m.List()
	if err != nil || len(list) == 0 {
		return nil, err
	}
	return &list[0], nil
}
