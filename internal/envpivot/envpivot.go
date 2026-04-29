// Package envpivot provides utilities for transposing secret maps
// by swapping keys and values, with optional collision handling.
package envpivot

import (
	"errors"
	"fmt"
	"sort"
)

// CollisionStrategy determines how duplicate values are handled during pivot.
type CollisionStrategy string

const (
	CollisionError  CollisionStrategy = "error"
	CollisionSkip   CollisionStrategy = "skip"
	CollisionSuffix CollisionStrategy = "suffix"
)

// Options configures the pivot operation.
type Options struct {
	Collision CollisionStrategy
}

// Pivoter swaps keys and values in a secrets map.
type Pivoter struct {
	opts Options
}

// New returns a Pivoter with the given options.
// If Collision is empty, CollisionError is used.
func New(opts Options) (*Pivoter, error) {
	if opts.Collision == "" {
		opts.Collision = CollisionError
	}
	switch opts.Collision {
	case CollisionError, CollisionSkip, CollisionSuffix:
	default:
		return nil, fmt.Errorf("envpivot: unknown collision strategy %q", opts.Collision)
	}
	return &Pivoter{opts: opts}, nil
}

// Apply returns a new map with keys and values swapped.
// Behaviour on duplicate values depends on the configured CollisionStrategy.
func (p *Pivoter) Apply(secrets map[string]string) (map[string]string, error) {
	result := make(map[string]string, len(secrets))
	// iterate in stable order for deterministic suffix numbering
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := secrets[k]
		if _, exists := result[v]; exists {
			switch p.opts.Collision {
			case CollisionError:
				return nil, fmt.Errorf("envpivot: duplicate value %q (from key %q)", v, k)
			case CollisionSkip:
				continue
			case CollisionSuffix:
				v = uniqueKey(result, v)
			}
		}
		result[v] = k
	}
	return result, nil
}

// Skipped returns keys whose values were skipped due to CollisionSkip.
func (p *Pivoter) Skipped(secrets map[string]string) []string {
	if p.opts.Collision != CollisionSkip {
		return nil
	}
	seen := map[string]bool{}
	var skipped []string
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := secrets[k]
		if seen[v] {
			skipped = append(skipped, k)
		}
		seen[v] = true
	}
	return skipped
}

func uniqueKey(m map[string]string, base string) string {
	for i := 2; ; i++ {
		candidate := fmt.Sprintf("%s_%d", base, i)
		if _, exists := m[candidate]; !exists {
			return candidate
		}
	}
}

var ErrEmptySecrets = errors.New("envpivot: secrets map is empty")
