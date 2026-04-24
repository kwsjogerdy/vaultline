package envdedupe_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envdedupe"
)

func TestFind_NoDuplicates(t *testing.T) {
	d := envdedupe.New(true)
	secrets := map[string]string{
		"A": "alpha",
		"B": "beta",
		"C": "gamma",
	}
	r := d.Find(secrets)
	if r.HasDuplicates() {
		t.Fatal("expected no duplicates")
	}
	if len(r.Removed) != 0 {
		t.Fatalf("expected no removals, got %v", r.Removed)
	}
	if len(r.Kept) != 3 {
		t.Fatalf("expected 3 kept, got %d", len(r.Kept))
	}
}

func TestFind_DetectsDuplicate(t *testing.T) {
	d := envdedupe.New(false)
	secrets := map[string]string{
		"X": "same",
		"Y": "same",
		"Z": "unique",
	}
	r := d.Find(secrets)
	if !r.HasDuplicates() {
		t.Fatal("expected duplicates")
	}
	if len(r.Duplicates) != 1 {
		t.Fatalf("expected 1 duplicate group, got %d", len(r.Duplicates))
	}
	dup := r.Duplicates[0]
	if dup.Value != "same" {
		t.Fatalf("unexpected value: %q", dup.Value)
	}
	if len(dup.Keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(dup.Keys))
	}
}

func TestFind_KeepFirst_RemovesLaterKeys(t *testing.T) {
	d := envdedupe.New(true)
	secrets := map[string]string{
		"AAA": "shared",
		"BBB": "shared",
		"CCC": "shared",
	}
	r := d.Find(secrets)
	if len(r.Removed) != 2 {
		t.Fatalf("expected 2 removed, got %v", r.Removed)
	}
	// AAA is alphabetically first — it should be kept
	if _, ok := r.Kept["AAA"]; !ok {
		t.Error("expected AAA to be kept")
	}
	for _, rem := range r.Removed {
		if _, ok := r.Kept[rem]; ok {
			t.Errorf("removed key %q should not appear in Kept", rem)
		}
	}
}

func TestFind_KeepFalse_RetainsAll(t *testing.T) {
	d := envdedupe.New(false)
	secrets := map[string]string{
		"P": "dup",
		"Q": "dup",
	}
	r := d.Find(secrets)
	if len(r.Removed) != 0 {
		t.Fatalf("expected no removals in report-only mode, got %v", r.Removed)
	}
	if len(r.Kept) != 2 {
		t.Fatalf("expected both keys kept, got %d", len(r.Kept))
	}
}

func TestFind_EmptySecrets(t *testing.T) {
	d := envdedupe.New(true)
	r := d.Find(map[string]string{})
	if r.HasDuplicates() {
		t.Fatal("empty map should have no duplicates")
	}
	if len(r.Kept) != 0 {
		t.Fatalf("expected empty kept map, got %d entries", len(r.Kept))
	}
}
