// Package envsampler provides random or deterministic sampling of secret maps,
// useful for previewing large secret sets without exposing all values.
package envsampler

import (
	"fmt"
	"math/rand"
	"sort"
)

// Options controls sampling behaviour.
type Options struct {
	// N is the maximum number of keys to return.
	N int
	// Seed makes sampling deterministic when non-zero.
	Seed int64
	// Prefix restricts sampling to keys with the given prefix.
	Prefix string
}

// Sampler draws a subset of secrets.
type Sampler struct {
	opts Options
	rng  *rand.Rand
}

// New returns a Sampler configured with opts.
// Returns an error if N < 1.
func New(opts Options) (*Sampler, error) {
	if opts.N < 1 {
		return nil, fmt.Errorf("envsampler: N must be >= 1, got %d", opts.N)
	}
	var src rand.Source
	if opts.Seed != 0 {
		src = rand.NewSource(opts.Seed)
	} else {
		src = rand.NewSource(rand.Int63())
	}
	return &Sampler{opts: opts, rng: rand.New(src)}, nil
}

// Sample returns up to N key-value pairs from secrets.
// If Prefix is set, only keys with that prefix are eligible.
// The returned map is a new map; the original is not modified.
func (s *Sampler) Sample(secrets map[string]string) map[string]string {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		if s.opts.Prefix != "" && len(k) < len(s.opts.Prefix) {
			continue
		}
		if s.opts.Prefix != "" && k[:len(s.opts.Prefix)] != s.opts.Prefix {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys) // deterministic base order before shuffle

	s.rng.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	if s.opts.N < len(keys) {
		keys = keys[:s.opts.N]
	}

	out := make(map[string]string, len(keys))
	for _, k := range keys {
		out[k] = secrets[k]
	}
	return out
}

// Keys returns only the sampled keys (sorted) without values.
func (s *Sampler) Keys(secrets map[string]string) []string {
	sampled := s.Sample(secrets)
	keys := make([]string, 0, len(sampled))
	for k := range sampled {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
