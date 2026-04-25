// Package envstats provides summary statistics over a secrets map.
package envstats

import (
	"fmt"
	"sort"
	"strings"
)

// Stats holds computed statistics for a secrets map.
type Stats struct {
	Total        int
	Empty        int
	Unique       int
	AvgValueLen  float64
	MaxValueLen  int
	MinValueLen  int
	PrefixCounts map[string]int
}

// Analyzer computes statistics over a secrets map.
type Analyzer struct {
	delimiter string
}

// New returns an Analyzer using the given prefix delimiter.
func New(delimiter string) *Analyzer {
	if delimiter == "" {
		delimiter = "_"
	}
	return &Analyzer{delimiter: delimiter}
}

// Compute calculates statistics for the given secrets map.
func (a *Analyzer) Compute(secrets map[string]string) Stats {
	if len(secrets) == 0 {
		return Stats{MinValueLen: 0}
	}

	seen := make(map[string]struct{})
	prefixCounts := make(map[string]int)
	totalLen := 0
	maxLen := 0
	minLen := -1
	empty := 0

	for k, v := range secrets {
		seen[v] = struct{}{}
		l := len(v)
		totalLen += l
		if l > maxLen {
			maxLen = l
		}
		if minLen == -1 || l < minLen {
			minLen = l
		}
		if v == "" {
			empty++
		}
		if idx := strings.Index(k, a.delimiter); idx > 0 {
			prefixCounts[k[:idx]]++
		} else {
			prefixCounts["(none)"]++
		}
	}

	if minLen == -1 {
		minLen = 0
	}

	avg := 0.0
	if len(secrets) > 0 {
		avg = float64(totalLen) / float64(len(secrets))
	}

	return Stats{
		Total:        len(secrets),
		Empty:        empty,
		Unique:       len(seen),
		AvgValueLen:  avg,
		MaxValueLen:  maxLen,
		MinValueLen:  minLen,
		PrefixCounts: prefixCounts,
	}
}

// Summary returns a human-readable summary of the stats.
func (s Stats) Summary() string {
	lines := []string{
		fmt.Sprintf("Total keys   : %d", s.Total),
		fmt.Sprintf("Empty values : %d", s.Empty),
		fmt.Sprintf("Unique values: %d", s.Unique),
		fmt.Sprintf("Avg val len  : %.1f", s.AvgValueLen),
		fmt.Sprintf("Max val len  : %d", s.MaxValueLen),
		fmt.Sprintf("Min val len  : %d", s.MinValueLen),
		"Prefix counts:",
	}

	prefixes := make([]string, 0, len(s.PrefixCounts))
	for p := range s.PrefixCounts {
		prefixes = append(prefixes, p)
	}
	sort.Strings(prefixes)
	for _, p := range prefixes {
		lines = append(lines, fmt.Sprintf("  %-20s %d", p, s.PrefixCounts[p]))
	}
	return strings.Join(lines, "\n")
}
