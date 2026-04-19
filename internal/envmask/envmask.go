// Package envmask provides selective masking of secrets based on key patterns
// before they are displayed or written to output.
package envmask

import (
	"regexp"
	"strings"
)

// Rule defines a masking rule for a key pattern.
type Rule struct {
	Pattern string
	Reveal  int // number of chars to reveal at end (0 = full redact)
}

// Masker applies masking rules to a map of secrets.
type Masker struct {
	rules   []compiledRule
	default_ string
}

type compiledRule struct {
	re     *regexp.Regexp
	reveal int
}

const defaultMask = "***"

// New creates a Masker with the given rules.
func New(rules []Rule) (*Masker, error) {
	compiled := make([]compiledRule, 0, len(rules))
	for _, r := range rules {
		re, err := regexp.Compile("(?i)" + r.Pattern)
		if err != nil {
			return nil, err
		}
		compiled = append(compiled, compiledRule{re: re, reveal: r.Reveal})
	}
	return &Masker{rules: compiled, default_: defaultMask}, nil
}

// Apply masks values in the secrets map according to the rules.
// Keys not matching any rule are returned as-is.
func (m *Masker) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = m.maskValue(k, v)
	}
	return out
}

// MaskValue masks a single value by key.
func (m *Masker) MaskValue(key, value string) string {
	return m.maskValue(key, value)
}

func (m *Masker) maskValue(key, value string) string {
	for _, r := range m.rules {
		if r.re.MatchString(key) {
			return applyMask(value, r.reveal)
		}
	}
	return value
}

func applyMask(value string, reveal int) string {
	if reveal <= 0 || reveal >= len(value) {
		return defaultMask
	}
	suffix := value[len(value)-reveal:]
	return strings.Repeat("*", 3) + suffix
}
