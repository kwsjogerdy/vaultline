// Package envlabel provides key-value labeling for secrets,
// allowing arbitrary metadata to be attached and queried.
package envlabel

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// Label is a metadata key-value pair attached to a secret key.
type Label struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Entry holds the labels for a single secret key.
type Entry struct {
	SecretKey string  `json:"secret_key"`
	Labels    []Label `json:"labels"`
}

// Labeler manages label metadata for secret keys.
type Labeler struct {
	path    string
	entries map[string]map[string]string // secretKey -> labelKey -> labelValue
}

// New creates a Labeler backed by the given JSON file path.
// If the file does not exist, an empty Labeler is returned.
func New(path string) (*Labeler, error) {
	l := &Labeler{
		path:    path,
		entries: make(map[string]map[string]string),
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return l, nil
	}
	if err != nil {
		return nil, fmt.Errorf("envlabel: read %s: %w", path, err)
	}
	var raw []Entry
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("envlabel: parse %s: %w", path, err)
	}
	for _, e := range raw {
		m := make(map[string]string, len(e.Labels))
		for _, lbl := range e.Labels {
			m[lbl.Key] = lbl.Value
		}
		l.entries[e.SecretKey] = m
	}
	return l, nil
}

// Set attaches a label key-value pair to a secret key.
func (l *Labeler) Set(secretKey, labelKey, labelValue string) {
	if _, ok := l.entries[secretKey]; !ok {
		l.entries[secretKey] = make(map[string]string)
	}
	l.entries[secretKey][labelKey] = labelValue
}

// Get returns the label value for a secret key and label key.
// The second return value is false if the label does not exist.
func (l *Labeler) Get(secretKey, labelKey string) (string, bool) {
	m, ok := l.entries[secretKey]
	if !ok {
		return "", false
	}
	v, ok := m[labelKey]
	return v, ok
}

// Filter returns secret keys whose label labelKey matches labelValue.
func (l *Labeler) Filter(labelKey, labelValue string) []string {
	var out []string
	for sk, m := range l.entries {
		if m[labelKey] == labelValue {
			out = append(out, sk)
		}
	}
	sort.Strings(out)
	return out
}

// Save persists the current label state to the backing file.
func (l *Labeler) Save() error {
	var raw []Entry
	keys := make([]string, 0, len(l.entries))
	for k := range l.entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, sk := range keys {
		e := Entry{SecretKey: sk}
		for lk, lv := range l.entries[sk] {
			e.Labels = append(e.Labels, Label{Key: lk, Value: lv})
		}
		sort.Slice(e.Labels, func(i, j int) bool { return e.Labels[i].Key < e.Labels[j].Key })
		raw = append(raw, e)
	}
	data, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return fmt.Errorf("envlabel: marshal: %w", err)
	}
	return os.WriteFile(l.path, data, 0o600)
}
