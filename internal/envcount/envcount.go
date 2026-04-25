// Package envcount provides utilities for counting and summarising
// secret entries by various dimensions such as prefix, value length,
// and emptiness.
package envcount

import (
	"fmt"
	"sort"
	"strings"
)

// Stats holds aggregated counts for a secrets map.
type Stats struct {
	Total    int
	Empty    int
	NonEmpty int
	// ByPrefix maps each detected top-level prefix (before the first '_') to
	// the number of keys that share it.
	ByPrefix map[string]int
}

// Summary returns a human-readable multi-line summary of the stats.
func (s Stats) Summary() string {
	lines := []string{
		fmt.Sprintf("total:     %d", s.Total),
		fmt.Sprintf("non-empty: %d", s.NonEmpty),
		fmt.Sprintf("empty:     %d", s.Empty),
	}
	if len(s.ByPrefix) > 0 {
		lines = append(lines, "by prefix:")
		prefixes := make([]string, 0, len(s.ByPrefix))
		for p := range s.ByPrefix {
			prefixes = append(prefixes, p)
		}
		sort.Strings(prefixes)
		for _, p := range prefixes {
			lines = append(lines, fmt.Sprintf("  %s: %d", p, s.ByPrefix[p]))
		}
	}
	return strings.Join(lines, "\n")
}

// Counter computes counts over a secrets map.
type Counter struct{}

// New returns a new Counter.
func New() *Counter { return &Counter{} }

// Compute analyses secrets and returns a Stats value.
func (c *Counter) Compute(secrets map[string]string) Stats {
	s := Stats{
		Total:    len(secrets),
		ByPrefix: make(map[string]int),
	}
	for k, v := range secrets {
		if strings.TrimSpace(v) == "" {
			s.Empty++
		} else {
			s.NonEmpty++
		}
		prefix := prefix(k)
		if prefix != "" {
			s.ByPrefix[prefix]++
		}
	}
	return s
}

// prefix returns the portion of key before the first '_', or empty string
// when no underscore is present.
func prefix(key string) string {
	if idx := strings.Index(key, "_"); idx > 0 {
		return key[:idx]
	}
	return ""
}
