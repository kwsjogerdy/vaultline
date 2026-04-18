// Package envchain supports chaining multiple secret sources with priority ordering.
package envchain

import "fmt"

// Source represents a named secret source.
type Source struct {
	Name     string
	Secrets  map[string]string
	Priority int // higher wins
}

// Chain holds an ordered list of sources.
type Chain struct {
	sources []Source
}

// New creates a new Chain.
func New(sources []Source) *Chain {
	sorted := make([]Source, len(sources))
	copy(sorted, sources)
	// simple insertion sort by priority descending
	for i := 1; i < len(sorted); i++ {
		for j := i; j > 0 && sorted[j].Priority > sorted[j-1].Priority; j-- {
			sorted[j], sorted[j-1] = sorted[j-1], sorted[j]
		}
	}
	return &Chain{sources: sorted}
}

// Resolve merges all sources, higher priority values win.
func (c *Chain) Resolve() map[string]string {
	result := make(map[string]string)
	// iterate lowest priority first so higher overwrites
	for i := len(c.sources) - 1; i >= 0; i-- {
		for k, v := range c.sources[i].Secrets {
			result[k] = v
		}
	}
	return result
}

// Origin returns the source name that provides the given key.
func (c *Chain) Origin(key string) (string, error) {
	for _, s := range c.sources {
		if _, ok := s.Secrets[key]; ok {
			return s.Name, nil
		}
	}
	return "", fmt.Errorf("envchain: key %q not found in any source", key)
}

// Sources returns the ordered source list.
func (c *Chain) Sources() []Source {
	return c.sources
}
