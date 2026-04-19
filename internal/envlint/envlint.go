// Package envlint provides linting rules for environment variable maps.
package envlint

import (
	"fmt"
	"regexp"
	"strings"
)

// Severity indicates how serious a lint finding is.
type Severity string

const (
	Warn  Severity = "WARN"
	Error Severity = "ERROR"
)

// Finding represents a single lint result.
type Finding struct {
	Key      string
	Rule     string
	Message  string
	Severity Severity
}

func (f Finding) String() string {
	return fmt.Sprintf("[%s] %s: %s (%s)", f.Severity, f.Key, f.Message, f.Rule)
}

// Linter runs a set of rules against a secrets map.
type Linter struct {
	allowEmpty bool
}

// New returns a Linter with default rules.
func New(allowEmpty bool) *Linter {
	return &Linter{allowEmpty: allowEmpty}
}

var validKeyRe = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// Lint checks secrets and returns all findings.
func (l *Linter) Lint(secrets map[string]string) []Finding {
	var findings []Finding

	for k, v := range secrets {
		if !validKeyRe.MatchString(k) {
			findings = append(findings, Finding{
				Key:      k,
				Rule:     "key-format",
				Message:  "key must be uppercase with underscores only",
				Severity: Error,
			})
		}

		if !l.allowEmpty && strings.TrimSpace(v) == "" {
			findings = append(findings, Finding{
				Key:      k,
				Rule:     "empty-value",
				Message:  "value is empty or whitespace",
				Severity: Warn,
			})
		}

		if strings.Contains(v, "\n") {
			findings = append(findings, Finding{
				Key:      k,
				Rule:     "multiline-value",
				Message:  "value contains newline characters",
				Severity: Warn,
			})
		}
	}

	return findings
}

// HasErrors returns true if any finding is Error severity.
func HasErrors(findings []Finding) bool {
	for _, f := range findings {
		if f.Severity == Error {
			return true
		}
	}
	return false
}
