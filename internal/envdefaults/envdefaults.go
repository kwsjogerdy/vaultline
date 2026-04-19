// Package envdefaults fills in missing secrets with predefined default values.
package envdefaults

import (
	"encoding/json"
	"errors"
	"os"
)

// Default represents a key with a fallback value and optional description.
type Default struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
}

// Applier holds a set of defaults to apply to a secrets map.
type Applier struct {
	defaults []Default
}

// New creates a new Applier with the given defaults.
func New(defaults []Default) (*Applier, error) {
	for _, d := range defaults {
		if d.Key == "" {
			return nil, errors.New("envdefaults: default key must not be empty")
		}
	}
	return &Applier{defaults: defaults}, nil
}

// Apply fills in any missing keys in secrets using the registered defaults.
func (a *Applier) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	for _, d := range a.defaults {
		if _, ok := out[d.Key]; !ok {
			out[d.Key] = d.Value
		}
	}
	return out
}

// Missing returns the list of default keys absent from secrets.
func (a *Applier) Missing(secrets map[string]string) []string {
	var missing []string
	for _, d := range a.defaults {
		if _, ok := secrets[d.Key]; !ok {
			missing = append(missing, d.Key)
		}
	}
	return missing
}

// LoadFile parses a JSON file of Default entries.
func LoadFile(path string) ([]Default, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var defaults []Default
	if err := json.Unmarshal(data, &defaults); err != nil {
		return nil, err
	}
	return defaults, nil
}
