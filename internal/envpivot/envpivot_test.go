package envpivot_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envpivot"
)

func TestApply_SwapsKeysAndValues(t *testing.T) {
	p, err := envpivot.New(envpivot.Options{})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	result, err := p.Apply(secrets)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if result["bar"] != "FOO" {
		t.Errorf("expected result[bar]=FOO, got %q", result["bar"])
	}
	if result["qux"] != "BAZ" {
		t.Errorf("expected result[qux]=BAZ, got %q", result["qux"])
	}
}

func TestApply_CollisionError(t *testing.T) {
	p, _ := envpivot.New(envpivot.Options{Collision: envpivot.CollisionError})
	secrets := map[string]string{"A": "same", "B": "same"}
	_, err := p.Apply(secrets)
	if err == nil {
		t.Fatal("expected error on duplicate value, got nil")
	}
}

func TestApply_CollisionSkip(t *testing.T) {
	p, _ := envpivot.New(envpivot.Options{Collision: envpivot.CollisionSkip})
	secrets := map[string]string{"A": "same", "B": "same"}
	result, err := p.Apply(secrets)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 entry, got %d", len(result))
	}
}

func TestApply_CollisionSuffix(t *testing.T) {
	p, _ := envpivot.New(envpivot.Options{Collision: envpivot.CollisionSuffix})
	secrets := map[string]string{"A": "dup", "B": "dup"}
	result, err := p.Apply(secrets)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 entries after suffix rename, got %d", len(result))
	}
	if _, ok := result["dup"]; !ok {
		t.Error("expected result[dup] to exist")
	}
	if _, ok := result["dup_2"]; !ok {
		t.Error("expected result[dup_2] to exist")
	}
}

func TestNew_UnknownStrategy_ReturnsError(t *testing.T) {
	_, err := envpivot.New(envpivot.Options{Collision: "bogus"})
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestSkipped_ReturnsSkippedKeys(t *testing.T) {
	p, _ := envpivot.New(envpivot.Options{Collision: envpivot.CollisionSkip})
	secrets := map[string]string{"A": "v", "B": "v", "C": "unique"}
	skipped := p.Skipped(secrets)
	if len(skipped) != 1 {
		t.Errorf("expected 1 skipped key, got %d", len(skipped))
	}
}

func TestApply_EmptyMap(t *testing.T) {
	p, _ := envpivot.New(envpivot.Options{})
	result, err := p.Apply(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}
