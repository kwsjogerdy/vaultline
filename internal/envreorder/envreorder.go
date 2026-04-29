// Package envreorder provides utilities for reordering secret map entries
// according to a user-defined key ordering or alphabetical fallback.
package envreorder

import (
	"errors"
	"sort"
)

// Reorderer reorders a map of secrets according to a defined key order.
type Reorderer struct {
	order    []string
	fallback bool // sort remaining keys alphabetically
}

// New creates a Reorderer with the given explicit key order.
// Keys not in the order list are placed at the end, sorted alphabetically
// if fallback is true, otherwise in their original iteration order.
func New(order []string, fallback bool) (*Reorderer, error) {
	if len(order) == 0 && !fallback {
		return nil, errors.New("envreorder: order list is empty and fallback is disabled")
	}
	return &Reorderer{order: order, fallback: fallback}, nil
}

// Apply returns a slice of key-value pairs in the defined order.
// Pairs whose keys appear in the order list come first (in order),
// followed by remaining keys.
func (r *Reorderer) Apply(secrets map[string]string) []KV {
	seen := make(map[string]bool, len(r.order))
	result := make([]KV, 0, len(secrets))

	for _, k := range r.order {
		if v, ok := secrets[k]; ok {
			result = append(result, KV{Key: k, Value: v})
			seen[k] = true
		}
	}

	remaining := make([]KV, 0)
	for k, v := range secrets {
		if !seen[k] {
			remaining = append(remaining, KV{Key: k, Value: v})
		}
	}

	if r.fallback {
		sort.Slice(remaining, func(i, j int) bool {
			return remaining[i].Key < remaining[j].Key
		})
	}

	return append(result, remaining...)
}

// Keys returns only the ordered key names from Apply.
func (r *Reorderer) Keys(secrets map[string]string) []string {
	kvs := r.Apply(secrets)
	keys := make([]string, len(kvs))
	for i, kv := range kvs {
		keys[i] = kv.Key
	}
	return keys
}

// KV holds a key-value pair preserving insertion order.
type KV struct {
	Key   string
	Value string
}
