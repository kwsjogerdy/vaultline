// Package envsplit splits a flat secrets map into multiple named buckets
// based on configurable key routing rules.
package envsplit

import (
	"errors"
	"sort"
	"strings"
)

// Rule defines how keys are routed to a named bucket.
type Rule struct {
	Bucket string
	Prefix string // keys with this prefix go to Bucket
}

// Result holds the split output.
type Result struct {
	Buckets  map[string]map[string]string
	Unrouted map[string]string
}

// Splitter routes secrets into named buckets.
type Splitter struct {
	rules []Rule
}

// New creates a Splitter with the provided rules.
// Returns an error if any rule has an empty bucket name.
func New(rules []Rule) (*Splitter, error) {
	for _, r := range rules {
		if strings.TrimSpace(r.Bucket) == "" {
			return nil, errors.New("envsplit: rule has empty bucket name")
		}
	}
	return &Splitter{rules: rules}, nil
}

// Apply routes each secret key into the first matching bucket.
// Keys that match no rule are placed in Unrouted.
func (s *Splitter) Apply(secrets map[string]string) Result {
	res := Result{
		Buckets:  make(map[string]map[string]string),
		Unrouted: make(map[string]string),
	}

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := secrets[k]
		routed := false
		for _, r := range s.rules {
			if r.Prefix == "" || strings.HasPrefix(k, r.Prefix) {
				if res.Buckets[r.Bucket] == nil {
					res.Buckets[r.Bucket] = make(map[string]string)
				}
				res.Buckets[r.Bucket][k] = v
				routed = true
				break
			}
		}
		if !routed {
			res.Unrouted[k] = v
		}
	}
	return res
}

// BucketNames returns a sorted list of bucket names produced by Apply.
func BucketNames(r Result) []string {
	names := make([]string, 0, len(r.Buckets))
	for n := range r.Buckets {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
