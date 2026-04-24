// Package envdependency resolves ordered secret loading based on declared dependencies.
package envdependency

import (
	"errors"
	"fmt"
)

// Rule declares that Key requires all Deps to be present.
type Rule struct {
	Key  string
	Deps []string
}

// Result holds the outcome of a dependency check for a single key.
type Result struct {
	Key     string
	Missing []string
	Satisfied bool
}

// Resolver checks secrets against declared dependency rules.
type Resolver struct {
	rules []Rule
}

// New creates a Resolver with the given rules.
func New(rules []Rule) (*Resolver, error) {
	for _, r := range rules {
		if r.Key == "" {
			return nil, errors.New("envdependency: rule key must not be empty")
		}
	}
	return &Resolver{rules: rules}, nil
}

// Check evaluates all rules against the provided secrets map.
// It returns one Result per rule.
func (r *Resolver) Check(secrets map[string]string) []Result {
	results := make([]Result, 0, len(r.rules))
	for _, rule := range r.rules {
		res := Result{Key: rule.Key}
		for _, dep := range rule.Deps {
			if _, ok := secrets[dep]; !ok {
				res.Missing = append(res.Missing, dep)
			}
		}
		res.Satisfied = len(res.Missing) == 0
		results = append(results, res)
	}
	return results
}

// Order returns keys in dependency-resolved order using a topological sort.
// Returns an error if a cycle is detected.
func Order(rules []Rule) ([]string, error) {
	graph := make(map[string][]string)
	nodes := make(map[string]struct{})
	for _, r := range rules {
		graph[r.Key] = r.Deps
		nodes[r.Key] = struct{}{}
		for _, d := range r.Deps {
			nodes[d] = struct{}{}
		}
	}

	visited := make(map[string]int) // 0=unvisited,1=visiting,2=done
	var order []string

	var visit func(n string) error
	visit = func(n string) error {
		switch visited[n] {
		case 2:
			return nil
		case 1:
			return fmt.Errorf("envdependency: cycle detected at %q", n)
		}
		visited[n] = 1
		for _, dep := range graph[n] {
			if err := visit(dep); err != nil {
				return err
			}
		}
		visited[n] = 2
		order = append(order, n)
		return nil
	}

	for n := range nodes {
		if err := visit(n); err != nil {
			return nil, err
		}
	}
	return order, nil
}
