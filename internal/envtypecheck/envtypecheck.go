// Package envtypecheck infers and validates the type of secret values.
package envtypecheck

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Type represents an inferred value type.
type Type string

const (
	TypeBool   Type = "bool"
	TypeInt    Type = "int"
	TypeFloat  Type = "float"
	TypeURL    Type = "url"
	TypeEmail  Type = "email"
	TypeString Type = "string"
)

var (
	urlRe   = regexp.MustCompile(`^https?://`)
	emailRe = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
)

// Result holds the inferred type and any validation error.
type Result struct {
	Key   string
	Value string
	Type  Type
	Error error
}

// Checker infers and optionally enforces types on secrets.
type Checker struct {
	enforce map[string]Type
}

// New creates a Checker. The enforce map pins specific keys to expected types.
func New(enforce map[string]Type) *Checker {
	if enforce == nil {
		enforce = map[string]Type{}
	}
	return &Checker{enforce: enforce}
}

// Infer returns the inferred Type for a single value.
func Infer(v string) Type {
	switch strings.ToLower(v) {
	case "true", "false", "yes", "no", "1", "0":
		return TypeBool
	}
	if _, err := strconv.Atoi(v); err == nil {
		return TypeInt
	}
	if _, err := strconv.ParseFloat(v, 64); err == nil {
		return TypeFloat
	}
	if urlRe.MatchString(v) {
		return TypeURL
	}
	if emailRe.MatchString(v) {
		return TypeEmail
	}
	return TypeString
}

// Check infers types for all secrets and validates against any enforced types.
func (c *Checker) Check(secrets map[string]string) []Result {
	results := make([]Result, 0, len(secrets))
	for k, v := range secrets {
		inferred := Infer(v)
		r := Result{Key: k, Value: v, Type: inferred}
		if expected, ok := c.enforce[k]; ok && inferred != expected {
			r.Error = fmt.Errorf("key %q: expected type %s, got %s", k, expected, inferred)
		}
		results = append(results, r)
	}
	return results
}

// HasErrors returns true if any result contains a type mismatch error.
func HasErrors(results []Result) bool {
	for _, r := range results {
		if r.Error != nil {
			return true
		}
	}
	return false
}
