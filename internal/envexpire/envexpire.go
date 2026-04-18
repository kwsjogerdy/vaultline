// Package envexpire tracks expiry timestamps for secrets and reports stale entries.
package envexpire

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// Entry holds expiry metadata for a single secret key.
type Entry struct {
	Key       string    `json:"key"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Store manages expiry entries persisted to a JSON file.
type Store struct {
	path    string
	entries map[string]Entry
}

// New loads (or creates) an expiry store at the given path.
func New(path string) (*Store, error) {
	s := &Store{path: path, entries: make(map[string]Entry)}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, fmt.Errorf("envexpire: read: %w", err)
	}
	var list []Entry
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("envexpire: parse: %w", err)
	}
	for _, e := range list {
		s.entries[e.Key] = e
	}
	return s, nil
}

// Set registers an expiry duration for a key.
func (s *Store) Set(key string, ttl time.Duration) {
	s.entries[key] = Entry{Key: key, ExpiresAt: time.Now().Add(ttl)}
}

// IsExpired reports whether the key exists and has passed its expiry time.
func (s *Store) IsExpired(key string) bool {
	e, ok := s.entries[key]
	if !ok {
		return false
	}
	return time.Now().After(e.ExpiresAt)
}

// Expired returns all keys that have passed their expiry time.
func (s *Store) Expired() []string {
	var out []string
	for k, e := range s.entries {
		if time.Now().After(e.ExpiresAt) {
			out = append(out, k)
		}
	}
	return out
}

// Save persists the store to disk.
func (s *Store) Save() error {
	list := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		list = append(list, e)
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return fmt.Errorf("envexpire: marshal: %w", err)
	}
	if err := os.WriteFile(s.path, data, 0600); err != nil {
		return fmt.Errorf("envexpire: write: %w", err)
	}
	return nil
}
