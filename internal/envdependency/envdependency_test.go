package envdependency_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envdependency"
)

func TestNew_EmptyKeyReturnsError(t *testing.T) {
	_, err := envdependency.New([]envdependency.Rule{{Key: "", Deps: []string{"A"}}})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestCheck_AllSatisfied(t *testing.T) {
	r, _ := envdependency.New([]envdependency.Rule{
		{Key: "DB_URL", Deps: []string{"DB_HOST", "DB_PORT"}},
	})
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	results := r.Check(secrets)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Satisfied {
		t.Errorf("expected satisfied, got missing: %v", results[0].Missing)
	}
}

func TestCheck_MissingDep(t *testing.T) {
	r, _ := envdependency.New([]envdependency.Rule{
		{Key: "APP_DSN", Deps: []string{"DB_HOST", "DB_PASS"}},
	})
	secrets := map[string]string{"DB_HOST": "localhost"}
	results := r.Check(secrets)
	if results[0].Satisfied {
		t.Error("expected unsatisfied")
	}
	if len(results[0].Missing) != 1 || results[0].Missing[0] != "DB_PASS" {
		t.Errorf("unexpected missing: %v", results[0].Missing)
	}
}

func TestCheck_NoDeps_AlwaysSatisfied(t *testing.T) {
	r, _ := envdependency.New([]envdependency.Rule{
		{Key: "STANDALONE", Deps: nil},
	})
	results := r.Check(map[string]string{})
	if !results[0].Satisfied {
		t.Error("expected satisfied when no deps")
	}
}

func TestOrder_SimpleChain(t *testing.T) {
	rules := []envdependency.Rule{
		{Key: "C", Deps: []string{"B"}},
		{Key: "B", Deps: []string{"A"}},
		{Key: "A", Deps: nil},
	}
	order, err := envdependency.Order(rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	pos := func(k string) int {
		for i, v := range order {
			if v == k {
				return i
			}
		}
		return -1
	}
	if pos("A") > pos("B") || pos("B") > pos("C") {
		t.Errorf("wrong order: %v", order)
	}
}

func TestOrder_CycleDetected(t *testing.T) {
	rules := []envdependency.Rule{
		{Key: "A", Deps: []string{"B"}},
		{Key: "B", Deps: []string{"A"}},
	}
	_, err := envdependency.Order(rules)
	if err == nil {
		t.Fatal("expected cycle error")
	}
}

func TestOrder_NoDeps_ReturnsAllKeys(t *testing.T) {
	rules := []envdependency.Rule{
		{Key: "X", Deps: nil},
		{Key: "Y", Deps: nil},
	}
	order, err := envdependency.Order(rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order) != 2 {
		t.Errorf("expected 2 keys, got %d", len(order))
	}
}
