package envdiff_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envdiff"
)

func TestCompare_AllMatch(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"A": "1", "B": "2"}
	r := envdiff.Compare(left, right)
	if r.HasDrift() {
		t.Fatal("expected no drift")
	}
}

func TestCompare_Changed(t *testing.T) {
	left := map[string]string{"A": "old"}
	right := map[string]string{"A": "new"}
	r := envdiff.Compare(left, right)
	if !r.HasDrift() {
		t.Fatal("expected drift")
	}
	if r.Entries[0].Status != "changed" {
		t.Fatalf("expected changed, got %s", r.Entries[0].Status)
	}
}

func TestCompare_LeftOnly(t *testing.T) {
	left := map[string]string{"X": "val"}
	right := map[string]string{}
	r := envdiff.Compare(left, right)
	if r.Entries[0].Status != "left_only" {
		t.Fatalf("expected left_only, got %s", r.Entries[0].Status)
	}
}

func TestCompare_RightOnly(t *testing.T) {
	left := map[string]string{}
	right := map[string]string{"Y": "val"}
	r := envdiff.Compare(left, right)
	if r.Entries[0].Status != "right_only" {
		t.Fatalf("expected right_only, got %s", r.Entries[0].Status)
	}
}

func TestSummary_Format(t *testing.T) {
	left := map[string]string{"A": "1", "B": "old", "C": "only"}
	right := map[string]string{"A": "1", "B": "new", "D": "only"}
	r := envdiff.Compare(left, right)
	s := r.Summary()
	if s == "" {
		t.Fatal("expected non-empty summary")
	}
	t.Log(s)
}

func TestCompare_Empty(t *testing.T) {
	r := envdiff.Compare(map[string]string{}, map[string]string{})
	if r.HasDrift() {
		t.Fatal("expected no drift for empty maps")
	}
}
