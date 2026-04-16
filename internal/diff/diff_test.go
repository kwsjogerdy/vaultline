package diff_test

import (
	"testing"

	"github.com/vaultline/internal/diff"
)

func TestCompare_Added(t *testing.T) {
	existing := map[string]string{}
	incoming := map[string]string{"NEW_KEY": "value"}
	r := diff.Compare(existing, incoming)
	if len(r.Added) != 1 || r.Added["NEW_KEY"] != "value" {
		t.Errorf("expected NEW_KEY in Added")
	}
}

func TestCompare_Removed(t *testing.T) {
	existing := map[string]string{"OLD_KEY": "value"}
	incoming := map[string]string{}
	r := diff.Compare(existing, incoming)
	if len(r.Removed) != 1 || r.Removed["OLD_KEY"] != "value" {
		t.Errorf("expected OLD_KEY in Removed")
	}
}

func TestCompare_Changed(t *testing.T) {
	existing := map[string]string{"KEY": "old"}
	incoming := map[string]string{"KEY": "new"}
	r := diff.Compare(existing, incoming)
	if len(r.Changed) != 1 || r.Changed["KEY"] != "new" {
		t.Errorf("expected KEY in Changed")
	}
}

func TestCompare_Unchanged(t *testing.T) {
	existing := map[string]string{"KEY": "same"}
	incoming := map[string]string{"KEY": "same"}
	r := diff.Compare(existing, incoming)
	if len(r.Unchanged) != 1 {
		t.Errorf("expected KEY in Unchanged")
	}
	if r.HasChanges() {
		t.Errorf("expected no changes")
	}
}

func TestHasChanges_False(t *testing.T) {
	r := diff.Compare(
		map[string]string{"A": "1"},
		map[string]string{"A": "1"},
	)
	if r.HasChanges() {
		t.Error("expected HasChanges to be false")
	}
}

func TestSummary(t *testing.T) {
	r := diff.Compare(
		map[string]string{"OLD": "x"},
		map[string]string{"NEW": "y"},
	)
	s := r.Summary()
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
