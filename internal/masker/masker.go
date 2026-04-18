// Package masker provides utilities for redacting sensitive secret values
// before displaying them in terminal output or audit logs.
package masker

import "strings"

const defaultMask = "****"

// Masker redacts secret values, optionally revealing a prefix.
type Masker struct {
	revealChars int
	mask        string
}

// New returns a Masker that fully redacts values.
func New() *Masker {
	return &Masker{mask: defaultMask}
}

// NewWithReveal returns a Masker that reveals the first n characters.
func NewWithReveal(n int) *Masker {
	return &Masker{revealChars: n, mask: defaultMask}
}

// Mask redacts a single value.
func (m *Masker) Mask(value string) string {
	if len(value) == 0 {
		return m.mask
	}
	if m.revealChars <= 0 || m.revealChars >= len(value) {
		return m.mask
	}
	return value[:m.revealChars] + m.mask
}

// MaskMap returns a copy of the map with all values redacted.
func (m *Masker) MaskMap(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = m.Mask(v)
	}
	return out
}

// MaskLine replaces secret values found in a line of text.
func (m *Masker) MaskLine(line string, secrets map[string]string) string {
	for _, v := range secrets {
		if v != "" {
			line = strings.ReplaceAll(line, v, m.mask)
		}
	}
	return line
}
