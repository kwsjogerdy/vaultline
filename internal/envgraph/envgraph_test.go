package envgraph_test

import (
	"strings"
	"testing"

	"github.com/vaultline/vaultline/internal/envgraph"
)

func TestNew_NoRefs_AllIsolated(t *testing.T) {
	secrets := map[string]string{
		"FOO": "bar",
		"BAZ": "qux",
	}
	g := envgraph.New(secrets)
	for _, n := range g.Nodes() {
		if len(n.Deps) != 0 {
			t.Errorf("expected no deps for %s, got %v", n.Key, n.Deps)
		}
	}
}

func TestNew_BuildsEdges(t *testing.T) {
	secrets := map[string]string{
		"BASE_URL": "https://example.com",
		"API_URL":  "${BASE_URL}/api",
		"FULL_URL": "${API_URL}/v1",
	}
	g := envgraph.New(secrets)
	nodes := map[string]*envgraph.Node{}
	for _, n := range g.Nodes() {
		nodes[n.Key] = n
	}
	if len(nodes["API_URL"].Deps) != 1 || nodes["API_URL"].Deps[0] != "BASE_URL" {
		t.Errorf("API_URL should depend on BASE_URL, got %v", nodes["API_URL"].Deps)
	}
	if len(nodes["FULL_URL"].Deps) != 1 || nodes["FULL_URL"].Deps[0] != "API_URL" {
		t.Errorf("FULL_URL should depend on API_URL, got %v", nodes["FULL_URL"].Deps)
	}
	if len(nodes["BASE_URL"].Deps) != 0 {
		t.Errorf("BASE_URL should have no deps, got %v", nodes["BASE_URL"].Deps)
	}
}

func TestNew_IgnoresUnknownRefs(t *testing.T) {
	secrets := map[string]string{
		"FOO": "${UNDEFINED_KEY}",
	}
	g := envgraph.New(secrets)
	for _, n := range g.Nodes() {
		if len(n.Deps) != 0 {
			t.Errorf("expected no deps for unknown ref, got %v", n.Deps)
		}
	}
}

func TestDOT_ContainsEdges(t *testing.T) {
	secrets := map[string]string{
		"A": "hello",
		"B": "${A}-world",
	}
	g := envgraph.New(secrets)
	dot := g.DOT()
	if !strings.Contains(dot, "digraph secrets") {
		t.Error("DOT output missing header")
	}
	if !strings.Contains(dot, `"B" -> "A"`) {
		t.Errorf("DOT missing edge B->A:\n%s", dot)
	}
}

func TestRoots_ReturnsUnreferenced(t *testing.T) {
	secrets := map[string]string{
		"BASE": "x",
		"MID":  "${BASE}",
		"TOP":  "${MID}",
	}
	g := envgraph.New(secrets)
	roots := g.Roots()
	if len(roots) != 1 || roots[0] != "TOP" {
		t.Errorf("expected [TOP] as root, got %v", roots)
	}
}

func TestRoots_EmptyGraph(t *testing.T) {
	g := envgraph.New(map[string]string{})
	if len(g.Roots()) != 0 {
		t.Error("expected no roots for empty graph")
	}
}
