// Package envprefix provides utilities for adding, removing, and replacing
// key prefixes in secret maps.
package envprefix

import "fmt"

// Transformer applies prefix operations to secret maps.
type Transformer struct {
	addPrefix    string
	removePrefix string
}

// New creates a new Transformer with the given add and remove prefixes.
// Either may be empty to skip that operation. Remove is applied before Add.
func New(addPrefix, removePrefix string) (*Transformer, error) {
	if addPrefix == "" && removePrefix == "" {
		return nil, fmt.Errorf("envprefix: at least one of addPrefix or removePrefix must be set")
	}
	return &Transformer{addPrefix: addPrefix, removePrefix: removePrefix}, nil
}

// Apply returns a new map with prefix transformations applied to all keys.
func (t *Transformer) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		newKey := k
		if t.removePrefix != "" {
			newKey = stripPrefix(newKey, t.removePrefix)
		}
		if t.addPrefix != "" {
			newKey = t.addPrefix + newKey
		}
		if newKey == "" {
			continue
		}
		out[newKey] = v
	}
	return out
}

// ReplacePrefix swaps oldPrefix for newPrefix on all matching keys.
func ReplacePrefix(secrets map[string]string, oldPrefix, newPrefix string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		newKey := stripPrefix(k, oldPrefix)
		if newKey != k {
			newKey = newPrefix + newKey
		}
		out[newKey] = v
	}
	return out
}

func stripPrefix(s, prefix string) string {
	if len(prefix) > 0 && len(s) >= len(prefix) && s[:len(prefix)] == prefix {
		return s[len(prefix):]
	}
	return s
}
