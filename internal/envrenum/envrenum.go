// Package envrenum validates that secret values belong to a defined set of allowed values.
package envrenum

import (
	"fmt"
	"sort"
	"strings"
)

// Rule defines an enumeration constraint for a single key.
type Rule struct {
	Key     string
	Allowed []string
}

// Violation describes a key whose value was not in the allowed set.
type Violation struct {
	Key     string
	Value   string
	Allowed []string
}

func (v Violation) Error() string {
	return fmt.Sprintf("key %q has value %q which is not in allowed set [%s]",
		v.Key, v.Value, strings.Join(v.Allowed, ", "))
}

// Validator checks secrets against enumeration rules.
type Validator struct {
	rules map[string][]string
}

// New creates a Validator from a slice of Rules.
func New(rules []Rule) (*Validator, error) {
	ruleMap := make(map[string][]string, len(rules))
	for _, r := range rules {
		if r.Key == "" {
			return nil, fmt.Errorf("rule has empty key")
		}
		if len(r.Allowed) == 0 {
			return nil, fmt.Errorf("rule for key %q has no allowed values", r.Key)
		}
		ruleMap[r.Key] = r.Allowed
	}
	return &Validator{rules: ruleMap}, nil
}

// Validate checks secrets against the registered rules and returns any violations.
func (v *Validator) Validate(secrets map[string]string) []Violation {
	var violations []Violation
	keys := make([]string, 0, len(v.rules))
	for k := range v.rules {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		allowed := v.rules[key]
		val, ok := secrets[key]
		if !ok {
			continue
		}
		if !contains(allowed, val) {
			violations = append(violations, Violation{Key: key, Value: val, Allowed: allowed})
		}
	}
	return violations
}

// HasViolations returns true if any violations exist.
func (v *Validator) HasViolations(secrets map[string]string) bool {
	return len(v.Validate(secrets)) > 0
}

func contains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}
