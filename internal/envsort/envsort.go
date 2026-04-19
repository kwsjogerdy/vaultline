// Package envsorт provides utilities for sorting secrets maps by key or value.
package envsort

import (
	"sort"
	"strings"
)

// Order defines the sort direction.
type Order string

const (
	Asc  Order = "asc"
	Desc Order = "desc"
)

// By defines the sort field.
type By string

const (
	ByKey   By = "key"
	ByValue By = "value"
)

// Options controls sorting behaviour.
type Options struct {
	By    By
	Order Order
}

// Sorter sorts secret maps.
type Sorter struct {
	opts Options
}

// New returns a Sorter with the given options.
// Defaults: sort by key ascending.
func New(opts Options) *Sorter {
	if opts.By == "" {
		opts.By = ByKey
	}
	if opts.Order == "" {
		opts.Order = Asc
	}
	return &Sorter{opts: opts}
}

// Sort returns a slice of [2]string pairs sorted according to the options.
func (s *Sorter) Sort(secrets map[string]string) [][2]string {
	pairs := make([][2]string, 0, len(secrets))
	for k, v := range secrets {
		pairs = append(pairs, [2]string{k, v})
	}

	sort.Slice(pairs, func(i, j int) bool {
		var a, b string
		if s.opts.By == ByValue {
			a, b = strings.ToLower(pairs[i][1]), strings.ToLower(pairs[j][1])
		} else {
			a, b = strings.ToLower(pairs[i][0]), strings.ToLower(pairs[j][0])
		}
		if s.opts.Order == Desc {
			return a > b
		}
		return a < b
	})

	return pairs
}

// Keys returns only the sorted keys.
func (s *Sorter) Keys(secrets map[string]string) []string {
	pairs := s.Sort(secrets)
	keys := make([]string, len(pairs))
	for i, p := range pairs {
		keys[i] = p[0]
	}
	return keys
}
