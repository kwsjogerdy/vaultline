// Package envtrim removes secrets whose values match a set of blank/placeholder patterns.
package envtrim

import "strings"

// DefaultPatterns are common placeholder values that should be trimmed.
var DefaultPatterns = []string{"", "null", "nil", "none", "undefined", "<none>", "<unset>", "todo", "fixme"}

// Trimmer removes placeholder secrets from a map.
type Trimmer struct {
	patterns map[string]struct{}
}

// New returns a Trimmer using the default placeholder patterns.
func New() *Trimmer {
	return NewWithPatterns(DefaultPatterns)
}

// NewWithPatterns returns a Trimmer using the supplied patterns (case-insensitive).
func NewWithPatterns(patterns []string) *Trimmer {
	m := make(map[string]struct{}, len(patterns))
	for _, p := range patterns {
		m[strings.ToLower(p)] = struct{}{}
	}
	return &Trimmer{patterns: m}
}

// Apply returns a copy of secrets with placeholder values removed.
func (t *Trimmer) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if _, skip := t.patterns[strings.ToLower(strings.TrimSpace(v))]; !skip {
			out[k] = v
		}
	}
	return out
}

// Removed returns keys that would be removed from secrets.
func (t *Trimmer) Removed(secrets map[string]string) []string {
	var keys []string
	for k, v := range secrets {
		if _, skip := t.patterns[strings.ToLower(strings.TrimSpace(v))]; skip {
			keys = append(keys, k)
		}
	}
	return keys
}
