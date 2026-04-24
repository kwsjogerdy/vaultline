// Package envindex builds a searchable index over a secrets map,
// supporting lookup by key substring, value substring, or label.
package envindex

import (
	"sort"
	"strings"
)

// Entry represents a single indexed secret.
type Entry struct {
	Key   string
	Value string
	Label string
}

// Result holds a matched entry and the field that matched.
type Result struct {
	Entry
	MatchedOn string // "key", "value", or "label"
}

// Index holds a pre-built search index.
type Index struct {
	entries []Entry
}

// New builds an Index from a secrets map and an optional label map.
// labelMap may be nil.
func New(secrets map[string]string, labelMap map[string]string) *Index {
	entries := make([]Entry, 0, len(secrets))
	for k, v := range secrets {
		lbl := ""
		if labelMap != nil {
			lbl = labelMap[k]
		}
		entries = append(entries, Entry{Key: k, Value: v, Label: lbl})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
	return &Index{entries: entries}
}

// Search returns all entries where the query appears in the key, value, or label.
// Matching is case-insensitive. Each entry appears at most once (first match wins).
func (idx *Index) Search(query string) []Result {
	q := strings.ToLower(query)
	var results []Result
	for _, e := range idx.entries {
		switch {
		case strings.Contains(strings.ToLower(e.Key), q):
			results = append(results, Result{Entry: e, MatchedOn: "key"})
		case strings.Contains(strings.ToLower(e.Value), q):
			results = append(results, Result{Entry: e, MatchedOn: "value"})
		case e.Label != "" && strings.Contains(strings.ToLower(e.Label), q):
			results = append(results, Result{Entry: e, MatchedOn: "label"})
		}
	}
	return results
}

// Keys returns all indexed keys in sorted order.
func (idx *Index) Keys() []string {
	keys := make([]string, len(idx.entries))
	for i, e := range idx.entries {
		keys[i] = e.Key
	}
	return keys
}

// Size returns the number of indexed entries.
func (idx *Index) Size() int {
	return len(idx.entries)
}
