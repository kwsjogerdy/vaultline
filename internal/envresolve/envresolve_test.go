package envresolve

import (
	"strings"
	"testing"
)

func TestResolve_NoReferences(t *testing.T) {
	r := New()
	in := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, err := r.Resolve(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestResolve_SimpleReference(t *testing.T) {
	r := New()
	in := map[string]string{
		"HOST": "localhost",
		"DSN":  "postgres://${HOST}/db",
	}
	out, err := r.Resolve(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DSN"] != "postgres://localhost/db" {
		t.Errorf("got %q", out["DSN"])
	}
}

func TestResolve_ChainedReferences(t *testing.T) {
	r := New()
	in := map[string]string{
		"A": "hello",
		"B": "${A}_world",
		"C": "${B}!",
	}
	out, err := r.Resolve(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["C"] != "hello_world!" {
		t.Errorf("got %q", out["C"])
	}
}

func TestResolve_UndefinedReference(t *testing.T) {
	r := New()
	in := map[string]string{"FOO": "${MISSING}"}
	_, err := r.Resolve(in)
	if err == nil {
		t.Fatal("expected error for undefined reference")
	}
	if !strings.Contains(err.Error(), "MISSING") {
		t.Errorf("error should mention key name, got: %v", err)
	}
}

func TestResolve_MaxDepthExceeded(t *testing.T) {
	r := &Resolver{maxDepth: 2}
	// A -> B -> A creates infinite recursion
	in := map[string]string{
		"A": "${B}",
		"B": "${A}",
	}
	_, err := r.Resolve(in)
	if err == nil {
		t.Fatal("expected error for max depth exceeded")
	}
}

func TestResolve_MultipleRefsInValue(t *testing.T) {
	r := New()
	in := map[string]string{
		"USER": "admin",
		"PASS": "secret",
		"URL":  "${USER}:${PASS}@host",
	}
	out, err := r.Resolve(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["URL"] != "admin:secret@host" {
		t.Errorf("got %q", out["URL"])
	}
}
