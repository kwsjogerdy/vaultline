// Package envsanitize provides utilities for sanitizing secret values
// by removing or replacing characters that may cause issues in shell
// environments, config parsers, or downstream consumers.
package envsanitize

import (
	"errors"
	"regexp"
	"strings"
)

// Rule describes a single sanitization rule applied to secret values.
type Rule struct {
	Pattern     string // regexp pattern to match
	Replacement string // replacement string (empty = remove)
}

// Options controls sanitizer behaviour.
type Options struct {
	Rules         []Rule
	StripNewlines bool
	StripNulls    bool
	TrimSpace     bool
}

// Sanitizer applies sanitization rules to a map of secrets.
type Sanitizer struct {
	opts     Options
	compiled []*regexp.Regexp
}

// DefaultOptions returns a sensible default set of sanitization options.
func DefaultOptions() Options {
	return Options{
		StripNewlines: true,
		StripNulls:    true,
		TrimSpace:     true,
	}
}

// New creates a new Sanitizer with the provided options.
func New(opts Options) (*Sanitizer, error) {
	s := &Sanitizer{opts: opts}
	for _, r := range opts.Rules {
		re, err := regexp.Compile(r.Pattern)
		if err != nil {
			return nil, errors.New("envsanitize: invalid pattern " + r.Pattern + ": " + err.Error())
		}
		s.compiled = append(s.compiled, re)
	}
	return s, nil
}

// Apply sanitizes all values in the provided secrets map and returns a new map.
func (s *Sanitizer) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = s.sanitize(v)
	}
	return out
}

// Sanitized returns the sanitized form of a single value.
func (s *Sanitizer) Sanitized(value string) string {
	return s.sanitize(value)
}

func (s *Sanitizer) sanitize(v string) string {
	if s.opts.StripNulls {
		v = strings.ReplaceAll(v, "\x00", "")
	}
	if s.opts.StripNewlines {
		v = strings.ReplaceAll(v, "\r\n", " ")
		v = strings.ReplaceAll(v, "\n", " ")
		v = strings.ReplaceAll(v, "\r", " ")
	}
	for i, re := range s.compiled {
		v = re.ReplaceAllString(v, s.opts.Rules[i].Replacement)
	}
	if s.opts.TrimSpace {
		v = strings.TrimSpace(v)
	}
	return v
}
