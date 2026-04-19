package envtrim_test

import (
	"sort"
	"testing"

	"github.com/vaultline/vaultline/internal/envtrim"
)

func TestApply_RemovesBlankValues(t *testing.T) {
	tr := envtrim.New()
	in := map[string]string{"A": "real", "B": "", "C": "value"}
	out := tr.Apply(in)
	if _, ok := out["B"]; ok {
		t.Fatal("expected B to be removed")
	}
	if out["A"] != "real" || out["C"] != "value" {
		t.Fatal("expected real values to be kept")
	}
}

func TestApply_RemovesPlaceholders(t *testing.T) {
	tr := envtrim.New()
	in := map[string]string{"X": "null", "Y": "NONE", "Z": "<unset>", "K": "keep"}
	out := tr.Apply(in)
	for _, key := range []string{"X", "Y", "Z"} {
		if _, ok := out[key]; ok {
			t.Fatalf("expected %s to be removed", key)
		}
	}
	if out["K"] != "keep" {
		t.Fatal("expected K to be kept")
	}
}

func TestApply_CaseInsensitive(t *testing.T) {
	tr := envtrim.New()
	in := map[string]string{"A": "NULL", "B": "Undefined"}
	out := tr.Apply(in)
	if len(out) != 0 {
		t.Fatalf("expected all removed, got %v", out)
	}
}

func TestApply_CustomPatterns(t *testing.T) {
	tr := envtrim.NewWithPatterns([]string{"placeholder", "tbd"})
	in := map[string]string{"A": "placeholder", "B": "real", "C": "TBD"}
	out := tr.Apply(in)
	if _, ok := out["A"]; ok {
		t.Fatal("expected A removed")
	}
	if _, ok := out["C"]; ok {
		t.Fatal("expected C removed (case-insensitive)")
	}
	if out["B"] != "real" {
		t.Fatal("expected B kept")
	}
}

func TestRemoved_ReturnsRemovedKeys(t *testing.T) {
	tr := envtrim.New()
	in := map[string]string{"A": "real", "B": "", "C": "null"}
	keys := tr.Removed(in)
	sort.Strings(keys)
	if len(keys) != 2 || keys[0] != "B" || keys[1] != "C" {
		t.Fatalf("unexpected removed keys: %v", keys)
	}
}

func TestApply_EmptyMap(t *testing.T) {
	tr := envtrim.New()
	out := tr.Apply(map[string]string{})
	if len(out) != 0 {
		t.Fatal("expected empty output")
	}
}
