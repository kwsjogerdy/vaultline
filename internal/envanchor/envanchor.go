// Package envanchor provides functionality for pinning specific secret keys
// to expected values and detecting when they drift from those anchored values.
// This is useful for catching accidental overwrites of critical configuration.
package envanchor

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"time"
)

// Anchor represents a pinned key-value expectation.
type Anchor struct {
	Key       string    `json:"key"`
	Expected  string    `json:"expected"`
	CreatedAt time.Time `json:"created_at"`
}

// Violation describes a key whose current value differs from its anchor.
type Violation struct {
	Key      string
	Expected string
	Actual   string
}

func (v Violation) String() string {
	return fmt.Sprintf("%s: expected %q, got %q", v.Key, v.Expected, v.Actual)
}

// Anchorer manages a set of anchored key expectations.
type Anchorer struct {
	anchors map[string]Anchor
	path    string
}

// New creates an Anchorer backed by the given file path.
// If the file does not exist, an empty set of anchors is returned.
func New(path string) (*Anchorer, error) {
	a := &Anchorer{
		anchors: make(map[string]Anchor),
		path:    path,
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return a, nil
	}
	if err != nil {
		return nil, fmt.Errorf("envanchor: read %s: %w", path, err)
	}
	var list []Anchor
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("envanchor: parse %s: %w", path, err)
	}
	for _, anc := range list {
		a.anchors[anc.Key] = anc
	}
	return a, nil
}

// Set adds or replaces an anchor for the given key using the provided expected value.
func (a *Anchorer) Set(key, expected string) error {
	if key == "" {
		return errors.New("envanchor: key must not be empty")
	}
	a.anchors[key] = Anchor{
		Key:       key,
		Expected:  expected,
		CreatedAt: time.Now().UTC(),
	}
	return a.save()
}

// Remove deletes the anchor for the given key, if present.
func (a *Anchorer) Remove(key string) error {
	delete(a.anchors, key)
	return a.save()
}

// Check compares the provided secrets map against all anchors and returns
// any violations where the current value differs from the expected value.
// Keys present in anchors but missing from secrets are also reported.
func (a *Anchorer) Check(secrets map[string]string) []Violation {
	var violations []Violation
	for key, anc := range a.anchors {
		actual, ok := secrets[key]
		if !ok {
			violations = append(violations, Violation{Key: key, Expected: anc.Expected, Actual: "<missing>"})
			continue
		}
		if actual != anc.Expected {
			violations = append(violations, Violation{Key: key, Expected: anc.Expected, Actual: actual})
		}
	}
	sort.Slice(violations, func(i, j int) bool {
		return violations[i].Key < violations[j].Key
	})
	return violations
}

// List returns all currently anchored keys in sorted order.
func (a *Anchorer) List() []Anchor {
	out := make([]Anchor, 0, len(a.anchors))
	for _, anc := range a.anchors {
		out = append(out, anc)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Key < out[j].Key
	})
	return out
}

func (a *Anchorer) save() error {
	list := a.List()
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return fmt.Errorf("envanchor: marshal: %w", err)
	}
	if err := os.WriteFile(a.path, data, 0o600); err != nil {
		return fmt.Errorf("envanchor: write %s: %w", a.path, err)
	}
	return nil
}
