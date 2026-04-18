// Package envmerge merges multiple secret maps with configurable conflict strategies.
package envmerge

import "fmt"

// Strategy defines how to handle key conflicts during merge.
type Strategy int

const (
	// StrategyFirst keeps the first value seen.
	StrategyFirst Strategy = iota
	// StrategyLast overwrites with the latest value.
	StrategyLast
	// StrategyError returns an error on conflict.
	StrategyError
)

// Result holds the merged secrets and metadata.
type Result struct {
	Secrets  map[string]string
	Conflicts []Conflict
}

// Conflict records a key that appeared in multiple sources.
type Conflict struct {
	Key    string
	Values []string
}

// Merger merges secret maps.
type Merger struct {
	strategy Strategy
}

// New creates a Merger with the given strategy.
func New(s Strategy) *Merger {
	return &Merger{strategy: s}
}

// Merge combines sources in order according to the configured strategy.
func (m *Merger) Merge(sources ...map[string]string) (*Result, error) {
	out := make(map[string]string)
	var conflicts []Conflict
	seen := make(map[string]string)

	for _, src := range sources {
		for k, v := range src {
			if existing, ok := seen[k]; ok {
				switch m.strategy {
				case StrategyError:
					return nil, fmt.Errorf("envmerge: conflict on key %q", k)
				case StrategyLast:
					out[k] = v
					conflicts = append(conflicts, Conflict{Key: k, Values: []string{existing, v}})
				case StrategyFirst:
					conflicts = append(conflicts, Conflict{Key: k, Values: []string{existing, v}})
				}
			} else {
				out[k] = v
				seen[k] = v
			}
		}
	}

	return &Result{Secrets: out, Conflicts: conflicts}, nil
}
