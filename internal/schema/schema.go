// Package schema validates that a set of secrets satisfies a required key schema.
package schema

import (
	"fmt"
	"strings"
)

// Rule describes a single key requirement.
type Rule struct {
	Key      string
	Required bool
	Pattern  string // optional: non-empty means value must contain this substring
}

// Schema holds a collection of rules.
type Schema struct {
	rules []Rule
}

// ValidationError holds all violations found during Validate.
type ValidationError struct {
	Violations []string
}

func (e *ValidationError) Error() string {
	return "schema validation failed:\n  " + strings.Join(e.Violations, "\n  ")
}

// New creates a Schema from a slice of Rules.
func New(rules []Rule) *Schema {
	return &Schema{rules: rules}
}

// Validate checks secrets against the schema rules.
// Returns a *ValidationError if any violations exist, nil otherwise.
func (s *Schema) Validate(secrets map[string]string) error {
	var violations []string
	for _, r := range s.rules {
		val, ok := secrets[r.Key]
		if r.Required && !ok {
			violations = append(violations, fmt.Sprintf("missing required key: %s", r.Key))
			continue
		}
		if ok && r.Pattern != "" && !strings.Contains(val, r.Pattern) {
			violations = append(violations, fmt.Sprintf("key %s: value does not contain expected pattern %q", r.Key, r.Pattern))
		}
	}
	if len(violations) > 0 {
		return &ValidationError{Violations: violations}
	}
	return nil
}
