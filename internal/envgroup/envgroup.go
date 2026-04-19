// Package envgroup groups secrets by logical namespace derived from key prefixes.
package envgroup

import (
	"fmt"
	"sort"
	"strings"
)

// Group holds secrets sharing a common prefix.
type Group struct {
	Name    string
	Secrets map[string]string
}

// Grouper splits a flat secret map into named groups.
type Grouper struct {
	sep string
}

// New returns a Grouper using sep as the prefix delimiter (e.g. "_").
func New(sep string) *Grouper {
	if sep == "" {
		sep = "_"
	}
	return &Grouper{sep: sep}
}

// Group partitions secrets by the first segment before sep.
// Keys with no delimiter land in the "default" group.
func (g *Grouper) Group(secrets map[string]string) []Group {
	buckets := map[string]map[string]string{}
	for k, v := range secrets {
		parts := strings.SplitN(k, g.sep, 2)
		var name string
		if len(parts) == 2 {
			name = strings.ToLower(parts[0])
		} else {
			name = "default"
		}
		if buckets[name] == nil {
			buckets[name] = map[string]string{}
		}
		buckets[name][k] = v
	}

	groups := make([]Group, 0, len(buckets))
	for name, m := range buckets {
		groups = append(groups, Group{Name: name, Secrets: m})
	}
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})
	return groups
}

// Summary returns a human-readable overview of each group.
func Summary(groups []Group) string {
	var sb strings.Builder
	for _, g := range groups {
		fmt.Fprintf(&sb, "[%s] %d key(s)\n", g.Name, len(g.Secrets))
	}
	return sb.String()
}
