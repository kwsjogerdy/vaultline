// Package profile manages named sync profiles for reusable configurations.
package profile

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Profile holds a named set of sync parameters.
type Profile struct {
	Name      string   `json:"name"`
	VaultPath string   `json:"vault_path"`
	OutputFile string  `json:"output_file"`
	Keys      []string `json:"keys,omitempty"`
	Exclude   []string `json:"exclude,omitempty"`
}

// Store manages profiles on disk.
type Store struct {
	dir string
}

// New returns a Store rooted at dir.
func New(dir string) *Store {
	return &Store{dir: dir}
}

// Save writes a profile to disk.
func (s *Store) Save(p Profile) error {
	if p.Name == "" {
		return errors.New("profile name must not be empty")
	}
	if err := os.MkdirAll(s.dir, 0700); err != nil {
		return err
	}
	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path(p.Name), b, 0600)
}

// Load reads a profile by name.
func (s *Store) Load(name string) (Profile, error) {
	b, err := os.ReadFile(s.path(name))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Profile{}, errors.New("profile not found: " + name)
		}
		return Profile{}, err
	}
	var p Profile
	return p, json.Unmarshal(b, &p)
}

// List returns all saved profile names.
func (s *Store) List() ([]string, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".json" {
			names = append(names, e.Name()[:len(e.Name())-5])
		}
	}
	return names, nil
}

// Delete removes a profile by name.
func (s *Store) Delete(name string) error {
	err := os.Remove(s.path(name))
	if errors.Is(err, os.ErrNotExist) {
		return errors.New("profile not found: " + name)
	}
	return err
}

func (s *Store) path(name string) string {
	return filepath.Join(s.dir, name+".json")
}
