// Package envfreeze allows marking secrets as frozen, preventing them from being overwritten during sync.
package envfreeze

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"time"
)

// Entry records when a key was frozen.
type Entry struct {
	Key       string    `json:"key"`
	FrozenAt  time.Time `json:"frozen_at"`
}

// Freezer manages a set of frozen secret keys.
type Freezer struct {
	path    string
	frozen  map[string]Entry
}

// New loads or initialises a Freezer backed by the given file path.
func New(path string) (*Freezer, error) {
	f := &Freezer{path: path, frozen: map[string]Entry{}}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return f, nil
	}
	if err != nil {
		return nil, err
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	for _, e := range entries {
		f.frozen[e.Key] = e
	}
	return f, nil
}

// Freeze marks a key as frozen.
func (f *Freezer) Freeze(key string) error {
	if key == "" {
		return errors.New("key must not be empty")
	}
	f.frozen[key] = Entry{Key: key, FrozenAt: time.Now().UTC()}
	return f.save()
}

// Unfreeze removes the freeze on a key.
func (f *Freezer) Unfreeze(key string) error {
	delete(f.frozen, key)
	return f.save()
}

// IsFrozen returns true if the key is currently frozen.
func (f *Freezer) IsFrozen(key string) bool {
	_, ok := f.frozen[key]
	return ok
}

// List returns all frozen entries sorted by key.
func (f *Freezer) List() []Entry {
	out := make([]Entry, 0, len(f.frozen))
	for _, e := range f.frozen {
		out = append(out, e)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out
}

// Filter removes frozen keys from the provided secrets map and returns the filtered copy.
func (f *Freezer) Filter(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if !f.IsFrozen(k) {
			out[k] = v
		}
	}
	return out
}

func (f *Freezer) save() error {
	entries := f.List()
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(f.path, data, 0o600)
}
