// Package envclean removes duplicate whitespace, normalizes line endings,
// and strips control characters from secret values.
package envclean

import (
	"strings"
	"unicode"
)

// Options controls which cleaning operations are applied.
type Options struct {
	TrimSpace        bool // trim leading/trailing whitespace from values
	CollapseSpaces   bool // collapse runs of internal whitespace to a single space
	StripControl     bool // remove non-printable control characters
	NormalizeNewlines bool // replace \r\n and bare \r with \n
}

// DefaultOptions returns a sensible default cleaning configuration.
func DefaultOptions() Options {
	return Options{
		TrimSpace:        true,
		CollapseSpaces:   false,
		StripControl:     true,
		NormalizeNewlines: true,
	}
}

// Cleaner applies cleaning rules to secret maps.
type Cleaner struct {
	opts    Options
	cleaned []string
}

// New creates a Cleaner with the given options.
func New(opts Options) *Cleaner {
	return &Cleaner{opts: opts}
}

// Apply cleans all values in the provided secrets map and returns a new map.
// The original map is not modified.
func (c *Cleaner) Apply(secrets map[string]string) map[string]string {
	c.cleaned = c.cleaned[:0]
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		cleaned := c.cleanValue(v)
		if cleaned != v {
			c.cleaned = append(c.cleaned, k)
		}
		out[k] = cleaned
	}
	return out
}

// Cleaned returns the keys whose values were modified by the last Apply call.
func (c *Cleaner) Cleaned() []string {
	result := make([]string, len(c.cleaned))
	copy(result, c.cleaned)
	return result
}

func (c *Cleaner) cleanValue(v string) string {
	if c.opts.NormalizeNewlines {
		v = strings.ReplaceAll(v, "\r\n", "\n")
		v = strings.ReplaceAll(v, "\r", "\n")
	}
	if c.opts.StripControl {
		v = strings.Map(func(r rune) rune {
			if r == '\n' || r == '\t' {
				return r
			}
			if unicode.IsControl(r) {
				return -1
			}
			return r
		}, v)
	}
	if c.opts.CollapseSpaces {
		fields := strings.Fields(v)
		v = strings.Join(fields, " ")
	}
	if c.opts.TrimSpace {
		v = strings.TrimSpace(v)
	}
	return v
}
