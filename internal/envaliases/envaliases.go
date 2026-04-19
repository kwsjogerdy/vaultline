// Package envaliases maps secret keys to one or more alias names.
package envaliases

import (
	"encoding/json"
	"fmt"
	"os"
)

// AliasMap maps an original key to a list of alias names.
type AliasMap map[string][]string

// Aliaser applies alias mappings to a secrets map.
type Aliaser struct {
	aliases AliasMap
}

// New creates an Aliaser with the given AliasMap.
func New(aliases AliasMap) *Aliaser {
	return &Aliaser{aliases: aliases}
}

// LoadFile reads an alias map from a JSON file.
func LoadFile(path string) (AliasMap, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("envaliases: read file: %w", err)
	}
	var m AliasMap
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("envaliases: parse json: %w", err)
	}
	return m, nil
}

// Apply returns a new map that includes the original secrets plus any
// aliased keys. If a source key does not exist the alias is skipped.
func (a *Aliaser) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	for src, names := range a.aliases {
		v, ok := secrets[src]
		if !ok {
			continue
		}
		for _, alias := range names {
			out[alias] = v
		}
	}
	return out
}

// Strip returns a copy of secrets with aliased keys removed, keeping
// only the original source keys.
func (a *Aliaser) Strip(secrets map[string]string) map[string]string {
	aliased := make(map[string]struct{})
	for _, names := range a.aliases {
		for _, alias := range names {
			aliased[alias] = struct{}{}
		}
	}
	out := make(map[string]string)
	for k, v := range secrets {
		if _, skip := aliased[k]; !skip {
			out[k] = v
		}
	}
	return out
}
