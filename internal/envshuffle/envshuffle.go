// Package envshuffle provides utilities for randomly shuffling or
// deterministically reordering secret maps for testing and sampling purposes.
package envshuffle

import (
	"fmt"
	"math/rand"
	"sort"
)

// Shuffler reorders secret map entries using a configurable strategy.
type Shuffler struct {
	seed    int64
	prefix  string
	reverse bool
}

// Option configures a Shuffler.
type Option func(*Shuffler)

// WithSeed sets a deterministic seed for reproducible shuffles.
func WithSeed(seed int64) Option {
	return func(s *Shuffler) { s.seed = seed }
}

// WithPrefix restricts shuffling to keys matching the given prefix.
func WithPrefix(prefix string) Option {
	return func(s *Shuffler) { s.prefix = prefix }
}

// WithReverse returns keys in reverse-sorted order instead of random.
func WithReverse() Option {
	return func(s *Shuffler) { s.reverse = true }
}

// New creates a Shuffler with the given options.
func New(opts ...Option) *Shuffler {
	s := &Shuffler{seed: 0}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Apply returns a new slice of [key, value] pairs in shuffled order.
// The original map is not mutated.
func (s *Shuffler) Apply(secrets map[string]string) ([]Entry, error) {
	if secrets == nil {
		return nil, fmt.Errorf("envshuffle: secrets map must not be nil")
	}

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		if s.prefix == "" || len(k) >= len(s.prefix) && k[:len(s.prefix)] == s.prefix {
			keys = append(keys, k)
		}
	}

	if s.reverse {
		sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	} else {
		r := rand.New(rand.NewSource(s.seed))
		r.Shuffle(len(keys), func(i, j int) { keys[i], keys[j] = keys[j], keys[i] })
	}

	out := make([]Entry, 0, len(keys))
	for _, k := range keys {
		out = append(out, Entry{Key: k, Value: secrets[k]})
	}
	return out, nil
}

// Entry holds a key-value pair from the shuffled output.
type Entry struct {
	Key   string
	Value string
}
