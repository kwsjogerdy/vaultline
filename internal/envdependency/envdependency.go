// Package envdependency resolves and validates declared dependencies
// between secret keys, ensuring that if a key is present its required
// peer keys are also present in the secret map.
package envdependency

import (
	"fmt"
	"sort"
	"strings"
)

// Rule declares that when Key is present, all keys listed in Requires
// must also be present in the secret map.
type Rule struct {
	Key      string   `json:"key"`
	Requires []string `json:"requires"`
}

// Violation describes a single unmet dependency.
type Violation struct {
	Key     string // the key that triggered the rule
	Missing string // the required key that was absent
}

func (v Violation) String() string {
	return fmt.Sprintf("%s requires %s (missing)", v.Key, v.Missing)
}

// Result holds the outcome of a dependency check.
type Result struct {
	Violations []Violation
}

// OK returns true when no violations were found.
func (r Result) OK() bool { return len(r.Violations) == 0 }

// Summary returns a human-readable multi-line summary.
func (r Result) Summary() string {
	if r.OK() {
		return "all dependencies satisfied"
	}
	lines := make([]string, 0, len(r.Violations))
	for _, v := range r.Violations {
		lines = append(lines, "  "+v.String())
	}
	sort.Strings(lines)
	return fmt.Sprintf("%d dependency violation(s):\n%s", len(r.Violations), strings.Join(lines, "\n"))
}

// Checker validates dependency rules against a secret map.
type Checker struct {
	rules []Rule
}

// New creates a Checker with the provided rules.
// Rules with an empty Key or no Requires entries are silently skipped.
func New(rules []Rule) *Checker {
	filtered := make([]Rule, 0, len(rules))
	for _, r := range rules {
		if r.Key == "" || len(r.Requires) == 0 {
			continue
		}
		filtered = append(filtered, r)
	}
	return &Checker{rules: filtered}
}

// Check evaluates all rules against secrets and returns a Result.
// Only rules whose Key exists in secrets are evaluated.
func (c *Checker) Check(secrets map[string]string) Result {
	var violations []Violation
	for _, rule := range c.rules {
		if _, present := secrets[rule.Key]; !present {
			continue
		}
		for _, req := range rule.Requires {
			if _, ok := secrets[req]; !ok {
				violations = append(violations, Violation{
					Key:     rule.Key,
					Missing: req,
				})
			}
		}
	}
	return Result{Violations: violations}
}
