// Package cache provides a simple file-based cache for Vault secrets
// to avoid redundant fetches within a TTL window.
package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Entry holds cached secrets and metadata.
type Entry struct {
	Secrets   map[string]string `json:"secrets"`
	FetchedAt time.Time         `json:"fetched_at"`
	Path      string            `json:"path"`
}

// Cache manages a local file-based secret cache.
type Cache struct {
	dir string
	ttl time.Duration
}

// New creates a Cache storing entries under dir with the given TTL.
func New(dir string, ttl time.Duration) *Cache {
	return &Cache{dir: dir, ttl: ttl}
}

func (c *Cache) filepath(vaultPath string) string {
	safe := filepath.Base(vaultPath)
	return filepath.Join(c.dir, safe+".cache.json")
}

// Get returns cached secrets if they exist and are within TTL.
func (c *Cache) Get(vaultPath string) (map[string]string, bool) {
	data, err := os.ReadFile(c.filepath(vaultPath))
	if err != nil {
		return nil, false
	}
	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, false
	}
	if time.Since(entry.FetchedAt) > c.ttl {
		return nil, false
	}
	return entry.Secrets, true
}

// Set writes secrets to the cache file.
func (c *Cache) Set(vaultPath string, secrets map[string]string) error {
	if err := os.MkdirAll(c.dir, 0700); err != nil {
		return err
	}
	entry := Entry{
		Secrets:   secrets,
		FetchedAt: time.Now(),
		Path:      vaultPath,
	}
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.filepath(vaultPath), data, 0600)
}

// Invalidate removes the cache entry for the given vault path.
func (c *Cache) Invalidate(vaultPath string) error {
	err := os.Remove(c.filepath(vaultPath))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
