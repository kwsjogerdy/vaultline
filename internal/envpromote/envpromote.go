// Package envpromote copies secrets from one environment profile to another,
// optionally transforming keys via a prefix swap.
package envpromote

import "fmt"

// Result holds the outcome of a promotion.
type Result struct {
	Copied  []string
	Skipped []string
}

// Promoter moves secrets between named environments.
type Promoter struct {
	src    string
	dst    string
	prefix string // strip from src keys before writing to dst
}

// New returns a Promoter from src to dst.
func New(src, dst string) *Promoter {
	return &Promoter{src: src, dst: dst}
}

// WithPrefix sets a key prefix to strip when promoting.
func (p *Promoter) WithPrefix(prefix string) *Promoter {
	p.prefix = prefix
	return p
}

// Promote copies secrets from src into dst, stripping the configured prefix.
// onlyKeys, if non-empty, limits which keys are promoted.
func (p *Promoter) Promote(src, dst map[string]string, onlyKeys []string) (map[string]string, Result) {
	allowed := toSet(onlyKeys)
	out := copyMap(dst)
	var res Result

	for k, v := range src {
		if len(allowed) > 0 && !allowed[k] {
			res.Skipped = append(res.Skipped, k)
			continue
		}
		newKey := stripPrefix(k, p.prefix)
		out[newKey] = v
		res.Copied = append(res.Copied, fmt.Sprintf("%s -> %s", k, newKey))
	}
	return out, res
}

func stripPrefix(key, prefix string) string {
	if prefix == "" || len(key) <= len(prefix) {
		return key
	}
	if key[:len(prefix)] == prefix {
		return key[len(prefix):]
	}
	return key
}

func toSet(keys []string) map[string]bool {
	m := make(map[string]bool, len(keys))
	for _, k := range keys {
		m[k] = true
	}
	return m
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
