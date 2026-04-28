// Package envpartition splits a secret map into two groups based on a predicate.
package envpartition

import "sort"

// Result holds the two partitions produced by Apply.
type Result struct {
	Matched    map[string]string
	Unmatched  map[string]string
}

// Partitioner splits secrets into matched / unmatched buckets.
type Partitioner struct {
	predicate func(key, value string) bool
}

// New creates a Partitioner using the supplied predicate.
// The predicate receives each key and value; returning true places the
// entry in the Matched bucket.
func New(predicate func(key, value string) bool) (*Partitioner, error) {
	if predicate == nil {
		return nil, errNilPredicate
	}
	return &Partitioner{predicate: predicate}, nil
}

// Apply partitions secrets and returns a Result.
func (p *Partitioner) Apply(secrets map[string]string) Result {
	matched := make(map[string]string)
	unmatched := make(map[string]string)
	for k, v := range secrets {
		if p.predicate(k, v) {
			matched[k] = v
		} else {
			unmatched[k] = v
		}
	}
	return Result{Matched: matched, Unmatched: unmatched}
}

// MatchedKeys returns a sorted slice of keys in the Matched bucket.
func (r Result) MatchedKeys() []string {
	return sortedKeys(r.Matched)
}

// UnmatchedKeys returns a sorted slice of keys in the Unmatched bucket.
func (r Result) UnmatchedKeys() []string {
	return sortedKeys(r.Unmatched)
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

var errNilPredicate = partitionError("predicate must not be nil")

type partitionError string

func (e partitionError) Error() string { return string(e) }
