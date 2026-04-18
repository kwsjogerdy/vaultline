// Package envvalidate checks env values against type and constraint rules.
package envvalidate

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Rule defines validation constraints for a single key.
type Rule struct {
	Required bool
	Type     string // "string", "int", "bool"
	Pattern  string // optional regex
}

// Result holds a single validation failure.
type Result struct {
	Key     string
	Message string
}

// Validator validates a map of secrets against a set of rules.
type Validator struct {
	rules map[string]Rule
}

// New creates a Validator with the given rules.
func New(rules map[string]Rule) *Validator {
	return &Validator{rules: rules}
}

// Validate checks secrets against all rules and returns any failures.
func (v *Validator) Validate(secrets map[string]string) []Result {
	var results []Result

	for key, rule := range v.rules {
		val, exists := secrets[key]

		if rule.Required && !exists {
			results = append(results, Result{Key: key, Message: "required key is missing"})
			continue
		}

		if !exists {
			continue
		}

		switch strings.ToLower(rule.Type) {
		case "int":
			if _, err := strconv.Atoi(val); err != nil {
				results = append(results, Result{Key: key, Message: fmt.Sprintf("expected int, got %q", val)})
			}
		case "bool":
			if _, err := strconv.ParseBool(val); err != nil {
				results = append(results, Result{Key: key, Message: fmt.Sprintf("expected bool, got %q", val)})
			}
		}

		if rule.Pattern != "" {
			re, err := regexp.Compile(rule.Pattern)
			if err != nil {
				results = append(results, Result{Key: key, Message: fmt.Sprintf("invalid pattern: %v", err)})
				continue
			}
			if !re.MatchString(val) {
				results = append(results, Result{Key: key, Message: fmt.Sprintf("value %q does not match pattern %q", val, rule.Pattern)})
			}
		}
	}

	return results
}
