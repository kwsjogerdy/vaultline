// Package envrewrite provides key/value rewriting rules for secret maps.
package envrewrite

import (
	"fmt"
	"strings"
)

// Rule defines a single rewrite operation.
type Rule struct {
	Find    string `json:"find"`
	Replace string `json:"replace"`
	Target  string `json:"target"` // "key" | "value" | "both"
}

// Rewriter applies a set of rewrite rules to a secret map.
type Rewriter struct {
	rules []Rule
}

// New creates a Rewriter with the given rules.
func New(rules []Rule) *Rewriter {
	return &Rewriter{rules: rules}
}

// Apply rewrites keys and/or values in secrets according to the configured rules.
func (r *Rewriter) Apply(secrets map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	for _, rule := range r.rules {
		if rule.Find == "" {
			return nil, fmt.Errorf("envrewrite: rule has empty find field")
		}
		target := rule.Target
		if target == "" {
			target = "key"
		}
		switch target {
		case "key":
			out = rewriteKeys(out, rule.Find, rule.Replace)
		case "value":
			out = rewriteValues(out, rule.Find, rule.Replace)
		case "both":
			out = rewriteKeys(out, rule.Find, rule.Replace)
			out = rewriteValues(out, rule.Find, rule.Replace)
		default:
			return nil, fmt.Errorf("envrewrite: unknown target %q (want key|value|both)", target)
		}
	}
	return out, nil
}

func rewriteKeys(m map[string]string, find, replace string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		newKey := strings.ReplaceAll(k, find, replace)
		out[newKey] = v
	}
	return out
}

func rewriteValues(m map[string]string, find, replace string) map[string]string {
	for k, v := range m {
		m[k] = strings.ReplaceAll(v, find, replace)
	}
	return m
}
