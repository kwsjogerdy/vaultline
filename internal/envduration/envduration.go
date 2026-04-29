// Package envduration parses and formats duration values stored in secret maps.
package envduration

import (
	"fmt"
	"sort"
	"time"
)

// Result holds the parsed duration alongside its original string value.
type Result struct {
	Key      string
	Raw      string
	Duration time.Duration
}

// Parser parses duration values from a secret map.
type Parser struct {
	keys []string // if non-empty, only these keys are parsed
}

// New returns a Parser. If keys is empty all entries are attempted.
func New(keys ...string) *Parser {
	return &Parser{keys: keys}
}

// Parse iterates over secrets and returns successfully parsed durations.
// Keys that cannot be parsed as durations are silently skipped unless they
// were explicitly requested, in which case an error is returned.
func (p *Parser) Parse(secrets map[string]string) ([]Result, error) {
	target := p.targetSet()

	var results []Result
	for _, k := range sortedKeys(secrets) {
		v := secrets[k]
		if len(target) > 0 && !target[k] {
			continue
		}
		d, err := time.ParseDuration(v)
		if err != nil {
			if target[k] {
				return nil, fmt.Errorf("envduration: key %q value %q: %w", k, v, err)
			}
			continue
		}
		results = append(results, Result{Key: k, Raw: v, Duration: d})
	}
	return results, nil
}

// Format converts duration values in seconds (as integers) back to Go duration
// strings and returns a new map with those keys replaced.
func Format(secrets map[string]string, keys ...string) map[string]string {
	out := make(map[string]string, len(secrets))
	set := make(map[string]bool, len(keys))
	for _, k := range keys {
		set[k] = true
	}
	for k, v := range secrets {
		if set[k] {
			if d, err := time.ParseDuration(v); err == nil {
				out[k] = d.String()
				continue
			}
		}
		out[k] = v
	}
	return out
}

func (p *Parser) targetSet() map[string]bool {
	if len(p.keys) == 0 {
		return nil
	}
	m := make(map[string]bool, len(p.keys))
	for _, k := range p.keys {
		m[k] = true
	}
	return m
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
