// Package envnormalize standardizes secret keys and values
// by applying configurable normalization rules.
package envnormalize

import (
	"strings"
	"unicode"
)

// Options controls normalization behaviour.
type Options struct {
	UppercaseKeys   bool
	TrimValues      bool
	ReplaceHyphens  bool // replace - with _ in keys
	StripNonAlnum   bool // strip non-alphanumeric/underscore chars from keys
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		UppercaseKeys:  true,
		TrimValues:     true,
		ReplaceHyphens: true,
		StripNonAlnum:  false,
	}
}

// Normalizer applies normalization rules to a secrets map.
type Normalizer struct {
	opts Options
}

// New creates a Normalizer with the given options.
func New(opts Options) *Normalizer {
	return &Normalizer{opts: opts}
}

// Apply returns a new map with normalized keys and values.
// Keys that become empty after normalization are dropped.
func (n *Normalizer) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		nk := n.normalizeKey(k)
		if nk == "" {
			continue
		}
		if n.opts.TrimValues {
			v = strings.TrimSpace(v)
		}
		out[nk] = v
	}
	return out
}

func (n *Normalizer) normalizeKey(k string) string {
	if n.opts.ReplaceHyphens {
		k = strings.ReplaceAll(k, "-", "_")
	}
	if n.opts.StripNonAlnum {
		var b strings.Builder
		for _, r := range k {
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
				b.WriteRune(r)
			}
		}
		k = b.String()
	}
	if n.opts.UppercaseKeys {
		k = strings.ToUpper(k)
	}
	return k
}
