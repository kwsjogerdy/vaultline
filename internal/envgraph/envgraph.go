// Package envgraph builds a dependency graph of secrets and renders
// it as an adjacency list or DOT format for visualization.
package envgraph

import (
	"fmt"
	"sort"
	"strings"
)

// Node represents a single secret key and its edges.
type Node struct {
	Key  string
	Deps []string // keys this node references
}

// Graph holds all nodes keyed by secret name.
type Graph struct {
	nodes map[string]*Node
}

// New constructs a Graph from a secrets map.
// Values that contain ${KEY} references are treated as edges.
func New(secrets map[string]string) *Graph {
	g := &Graph{nodes: make(map[string]*Node, len(secrets))}
	for k := range secrets {
		g.nodes[k] = &Node{Key: k}
	}
	for k, v := range secrets {
		deps := extractRefs(v)
		for _, dep := range deps {
			if _, ok := g.nodes[dep]; ok {
				g.nodes[k].Deps = append(g.nodes[k].Deps, dep)
			}
		}
		sort.Strings(g.nodes[k].Deps)
	}
	return g
}

// Nodes returns all nodes sorted by key.
func (g *Graph) Nodes() []*Node {
	out := make([]*Node, 0, len(g.nodes))
	for _, n := range g.nodes {
		out = append(out, n)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out
}

// DOT renders the graph in Graphviz DOT format.
func (g *Graph) DOT() string {
	var sb strings.Builder
	sb.WriteString("digraph secrets {\n")
	for _, n := range g.Nodes() {
		for _, dep := range n.Deps {
			fmt.Fprintf(&sb, "  %q -> %q;\n", n.Key, dep)
		}
		if len(n.Deps) == 0 {
			fmt.Fprintf(&sb, "  %q;\n", n.Key)
		}
	}
	sb.WriteString("}\n")
	return sb.String()
}

// Roots returns keys that no other node depends on.
func (g *Graph) Roots() []string {
	depended := make(map[string]bool)
	for _, n := range g.nodes {
		for _, d := range n.Deps {
			depended[d] = true
		}
	}
	var roots []string
	for k := range g.nodes {
		if !depended[k] {
			roots = append(roots, k)
		}
	}
	sort.Strings(roots)
	return roots
}

// extractRefs returns all ${VAR} references found in s.
func extractRefs(s string) []string {
	var refs []string
	for {
		start := strings.Index(s, "${")
		if start == -1 {
			break
		}
		end := strings.Index(s[start:], "}")
		if end == -1 {
			break
		}
		refs = append(refs, s[start+2:start+end])
		s = s[start+end+1:]
	}
	return refs
}
