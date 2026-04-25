// Package envwildcard provides glob-style wildcard matching and filtering
// for secret maps, supporting patterns like "DB_*", "*_SECRET", and "APP_*_KEY".
package envwildcard

import (
	"fmt"
	"path"
	"sort"
)

// Matcher filters secrets using glob patterns.
type Matcher struct {
	include []string
	exclude []string
}

// New creates a Matcher with the given include and exclude glob patterns.
// Returns an error if any pattern is syntactically invalid.
func New(include, exclude []string) (*Matcher, error) {
	for _, p := range append(include, exclude...) {
		if _, err := path.Match(p, ""); err != nil {
			return nil, fmt.Errorf("envwildcard: invalid pattern %q: %w", p, err)
		}
	}
	return &Matcher{include: include, exclude: exclude}, nil
}

// Apply returns a filtered copy of secrets where keys match at least one
// include pattern (or all keys if no include patterns are set) and do not
// match any exclude pattern.
func (m *Matcher) Apply(secrets map[string]string) (map[string]string, error) {
	out := make(map[string]string)
	for k, v := range secrets {
		if len(m.include) > 0 {
			matched, err := matchAny(k, m.include)
			if err != nil {
				return nil, err
			}
			if !matched {
				continue
			}
		}
		blocked, err := matchAny(k, m.exclude)
		if err != nil {
			return nil, err
		}
		if blocked {
			continue
		}
		out[k] = v
	}
	return out, nil
}

// Matched returns the sorted list of keys from secrets that would survive Apply.
func (m *Matcher) Matched(secrets map[string]string) ([]string, error) {
	filtered, err := m.Apply(secrets)
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(filtered))
	for k := range filtered {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys, nil
}

func matchAny(key string, patterns []string) (bool, error) {
	for _, p := range patterns {
		ok, err := path.Match(p, key)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}
