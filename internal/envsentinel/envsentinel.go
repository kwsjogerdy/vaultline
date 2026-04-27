// Package envsentinel watches for secrets that match alert conditions
// and reports violations (e.g. keys with empty values, forbidden patterns).
package envsentinel

import (
	"fmt"
	"regexp"
	"sort"
)

// Rule defines a single sentinel check.
type Rule struct {
	// Name is a human-readable label for the rule.
	Name string
	// Pattern is a regex matched against the key name.
	Pattern string
	// ForbidEmpty triggers a violation when the matched key has an empty value.
	ForbidEmpty bool
	// ForbidPattern is a regex that must NOT match the value.
	ForbidPattern string
}

// Violation describes a rule breach for a specific key.
type Violation struct {
	Rule    string
	Key     string
	Message string
}

func (v Violation) String() string {
	return fmt.Sprintf("[%s] %s: %s", v.Rule, v.Key, v.Message)
}

// Sentinel evaluates a set of rules against a secrets map.
type Sentinel struct {
	rules []Rule
}

// New creates a Sentinel with the provided rules.
func New(rules []Rule) (*Sentinel, error) {
	for _, r := range rules {
		if _, err := regexp.Compile(r.Pattern); err != nil {
			return nil, fmt.Errorf("rule %q: invalid pattern: %w", r.Name, err)
		}
		if r.ForbidPattern != "" {
			if _, err := regexp.Compile(r.ForbidPattern); err != nil {
				return nil, fmt.Errorf("rule %q: invalid forbid_pattern: %w", r.Name, err)
			}
		}
	}
	return &Sentinel{rules: rules}, nil
}

// Check evaluates all rules against secrets and returns any violations.
func (s *Sentinel) Check(secrets map[string]string) []Violation {
	var violations []Violation

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, rule := range s.rules {
		keyRe := regexp.MustCompile(rule.Pattern)
		var forbidRe *regexp.Regexp
		if rule.ForbidPattern != "" {
			forbidRe = regexp.MustCompile(rule.ForbidPattern)
		}

		for _, k := range keys {
			if !keyRe.MatchString(k) {
				continue
			}
			v := secrets[k]
			if rule.ForbidEmpty && v == "" {
				violations = append(violations, Violation{
					Rule:    rule.Name,
					Key:     k,
					Message: "value must not be empty",
				})
			}
			if forbidRe != nil && forbidRe.MatchString(v) {
				violations = append(violations, Violation{
					Rule:    rule.Name,
					Key:     k,
					Message: fmt.Sprintf("value matches forbidden pattern %q", rule.ForbidPattern),
				})
			}
		}
	}
	return violations
}

// HasViolations returns true when Check produces at least one violation.
func (s *Sentinel) HasViolations(secrets map[string]string) bool {
	return len(s.Check(secrets)) > 0
}
