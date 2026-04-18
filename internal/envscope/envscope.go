// Package envscope provides environment scoping — restricting which secrets
// are visible based on a named scope (e.g. "production", "staging").
package envscope

import (
	"fmt"
	"strings"
)

// Scope represents a named environment scope with allowed key prefixes.
type Scope struct {
	Name     string
	Prefixes []string
}

// Manager manages multiple scopes.
type Manager struct {
	scopes map[string]Scope
}

// New returns a new Manager with the given scopes.
func New(scopes []Scope) *Manager {
	m := &Manager{scopes: make(map[string]Scope, len(scopes))}
	for _, s := range scopes {
		m.scopes[strings.ToLower(s.Name)] = s
	}
	return m
}

// Apply filters secrets to only those allowed by the named scope.
// If the scope has no prefixes, all keys are allowed.
func (m *Manager) Apply(scope string, secrets map[string]string) (map[string]string, error) {
	s, ok := m.scopes[strings.ToLower(scope)]
	if !ok {
		return nil, fmt.Errorf("envscope: unknown scope %q", scope)
	}
	if len(s.Prefixes) == 0 {
		result := make(map[string]string, len(secrets))
		for k, v := range secrets {
			result[k] = v
		}
		return result, nil
	}
	result := make(map[string]string)
	for k, v := range secrets {
		for _, p := range s.Prefixes {
			if strings.HasPrefix(k, p) {
				result[k] = v
				break
			}
		}
	}
	return result, nil
}

// List returns all registered scope names.
func (m *Manager) List() []string {
	names := make([]string, 0, len(m.scopes))
	for name := range m.scopes {
		names = append(names, name)
	}
	return names
}
