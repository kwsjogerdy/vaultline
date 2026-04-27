// Package envquote normalises quoting styles across secret values.
// It can enforce single-quotes, double-quotes, or strip all quotes.
package envquote

import (
	"fmt"
	"strings"
)

// Style describes the quoting style to apply.
type Style string

const (
	StyleDouble Style = "double"
	StyleSingle Style = "single"
	StyleStrip  Style = "strip"
)

// Quoter applies a quoting style to secret values.
type Quoter struct {
	style Style
}

// New returns a Quoter for the given style.
// Returns an error if the style is not recognised.
func New(style Style) (*Quoter, error) {
	switch style {
	case StyleDouble, StyleSingle, StyleStrip:
		return &Quoter{style: style}, nil
	}
	return nil, fmt.Errorf("envquote: unknown style %q (want double|single|strip)", style)
}

// Apply returns a new map with values re-quoted according to the style.
func (q *Quoter) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = q.transform(v)
	}
	return out
}

// transform strips existing outer quotes then re-applies the target style.
func (q *Quoter) transform(v string) string {
	bare := stripOuter(v)
	switch q.style {
	case StyleDouble:
		return `"` + bare + `"`
	case StyleSingle:
		return `'` + bare + `'`
	default: // StyleStrip
		return bare
	}
}

// stripOuter removes a single layer of matching outer quotes if present.
func stripOuter(v string) string {
	if len(v) >= 2 {
		if (strings.HasPrefix(v, `"`) && strings.HasSuffix(v, `"`)) ||
			(strings.HasPrefix(v, `'`) && strings.HasSuffix(v, `'`)) {
			return v[1 : len(v)-1]
		}
	}
	return v
}
