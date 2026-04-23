// Package envredact provides selective redaction of secret values
// based on configurable key patterns before output or logging.
package envredact

import (
	"regexp"
	"strings"
)

// Result holds the redacted map and which keys were redacted.
type Result struct {
	Secrets  map[string]string
	Redacted []string
}

// Redactor applies redaction rules to a secrets map.
type Redactor struct {
	patterns []*regexp.Regexp
	placeholder string
}

var defaultPatterns = []string{
	`(?i)password`,
	`(?i)secret`,
	`(?i)token`,
	`(?i)api[_-]?key`,
	`(?i)private[_-]?key`,
	`(?i)credentials`,
	`(?i)auth`,
}

// New creates a Redactor using the default sensitive key patterns.
func New() *Redactor {
	return NewWithPatterns(defaultPatterns, "[REDACTED]")
}

// NewWithPatterns creates a Redactor with custom patterns and placeholder.
func NewWithPatterns(patterns []string, placeholder string) *Redactor {
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		if re, err := regexp.Compile(p); err == nil {
			compiled = append(compiled, re)
		}
	}
	if placeholder == "" {
		placeholder = "[REDACTED]"
	}
	return &Redactor{patterns: compiled, placeholder: placeholder}
}

// Apply returns a copy of secrets with sensitive values replaced by the placeholder.
func (r *Redactor) Apply(secrets map[string]string) Result {
	out := make(map[string]string, len(secrets))
	var redacted []string
	for k, v := range secrets {
		if r.isSensitive(k) {
			out[k] = r.placeholder
			redacted = append(redacted, k)
		} else {
			out[k] = v
		}
	}
	return Result{Secrets: out, Redacted: redacted}
}

// IsSensitive reports whether a key matches any redaction pattern.
func (r *Redactor) IsSensitive(key string) bool {
	return r.isSensitive(key)
}

func (r *Redactor) isSensitive(key string) bool {
	upper := strings.ToUpper(key)
	_ = upper
	for _, re := range r.patterns {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}
