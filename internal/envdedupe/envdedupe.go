// Package envdedupe detects and removes duplicate values across secret maps.
package envdedupe

import "sort"

// Duplicate holds keys that share the same value.
type Duplicate struct {
	Value string
	Keys  []string
}

// Result is returned by Find.
type Result struct {
	Duplicates []Duplicate
	Kept       map[string]string
	Removed    []string
}

// Deduper finds and optionally removes keys with duplicate values.
type Deduper struct {
	keepFirst bool
}

// New returns a Deduper. When keepFirst is true, the first
// (alphabetically) key for each value is kept; others are removed.
func New(keepFirst bool) *Deduper {
	return &Deduper{keepFirst: keepFirst}
}

// Find scans secrets for duplicate values and returns a Result.
func (d *Deduper) Find(secrets map[string]string) Result {
	// group keys by value
	index := make(map[string][]string)
	for k, v := range secrets {
		index[v] = append(index[v], k)
	}

	var dups []Duplicate
	removed := []string{}
	kept := make(map[string]string, len(secrets))

	for v, keys := range index {
		if len(keys) < 2 {
			// no duplicate — always keep
			kept[keys[0]] = v
			continue
		}

		sort.Strings(keys)
		dups = append(dups, Duplicate{Value: v, Keys: keys})

		if d.keepFirst {
			kept[keys[0]] = v
			removed = append(removed, keys[1:]...)
		} else {
			for _, k := range keys {
				kept[k] = v
			}
		}
	}

	// stable sort for deterministic output
	sort.Slice(dups, func(i, j int) bool {
		return dups[i].Keys[0] < dups[j].Keys[0]
	})
	sort.Strings(removed)

	return Result{
		Duplicates: dups,
		Kept:       kept,
		Removed:    removed,
	}
}

// HasDuplicates returns true when at least one duplicate value was found.
func (r Result) HasDuplicates() bool {
	return len(r.Duplicates) > 0
}
