// Package redact provides utilities for scrubbing sensitive keys from output.
package redact

import "strings"

// DefaultSensitivePatterns contains common substrings that indicate a secret.
var DefaultSensitivePatterns = []string{
	"password", "passwd", "secret", "token", "api_key", "apikey",
	"private_key", "auth", "credential", "cert",
}

// Redactor scrubs values whose keys match sensitive patterns.
type Redactor struct {
	patterns []string
	placeholder string
}

// New returns a Redactor using DefaultSensitivePatterns.
func New() *Redactor {
	return NewWithPatterns(DefaultSensitivePatterns, "[REDACTED]")
}

// NewWithPatterns returns a Redactor with custom patterns and placeholder.
func NewWithPatterns(patterns []string, placeholder string) *Redactor {
	norm := make([]string, len(patterns))
	for i, p := range patterns {
		norm[i] = strings.ToLower(p)
	}
	return &Redactor{patterns: norm, placeholder: placeholder}
}

// IsSensitive reports whether the given key matches any sensitive pattern.
func (r *Redactor) IsSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, p := range r.patterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

// RedactMap returns a copy of m with sensitive values replaced by the placeholder.
func (r *Redactor) RedactMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		if r.IsSensitive(k) {
			out[k] = r.placeholder
		} else {
			out[k] = v
		}
	}
	return out
}

// RedactValue returns the placeholder if key is sensitive, otherwise value.
func (r *Redactor) RedactValue(key, value string) string {
	if r.IsSensitive(key) {
		return r.placeholder
	}
	return value
}
