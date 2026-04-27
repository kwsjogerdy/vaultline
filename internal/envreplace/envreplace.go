// Package envreplace provides find-and-replace operations on secret values.
package envreplace

import (
	"errors"
	"strings"
)

// Rule defines a single find-and-replace operation.
type Rule struct {
	Find    string
	Replace string
	Keys    []string // if empty, applies to all keys
}

// Replacer applies substitution rules to secret maps.
type Replacer struct {
	rules []Rule
}

// New creates a Replacer from the given rules.
// Returns an error if any rule has an empty Find string.
func New(rules []Rule) (*Replacer, error) {
	for _, r := range rules {
		if r.Find == "" {
			return nil, errors.New("envreplace: find string must not be empty")
		}
	}
	return &Replacer{rules: rules}, nil
}

// Apply returns a new map with all rules applied to the values.
func (r *Replacer) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	for _, rule := range r.rules {
		if len(rule.Keys) == 0 {
			for k, v := range out {
				out[k] = strings.ReplaceAll(v, rule.Find, rule.Replace)
			}
		} else {
			for _, k := range rule.Keys {
				if v, ok := out[k]; ok {
					out[k] = strings.ReplaceAll(v, rule.Find, rule.Replace)
				}
			}
		}
	}
	return out
}

// Changed returns the keys whose values differ after Apply.
func (r *Replacer) Changed(secrets map[string]string) []string {
	applied := r.Apply(secrets)
	var keys []string
	for k, v := range secrets {
		if applied[k] != v {
			keys = append(keys, k)
		}
	}
	return keys
}
