package envcompare_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envcompare"
)

func TestCompare_AllMatch(t *testing.T) {
	c := envcompare.New("dev", "staging")
	r := c.Compare(map[string]string{"A": "1"}, map[string]string{"A": "1"})
	if len(r.Results) != 1 || r.Results[0].Status != "match" {
		t.Fatalf("expected match, got %+v", r.Results)
	}
}

func TestCompare_Differ(t *testing.T) {
	c := envcompare.New("dev", "staging")
	r := c.Compare(map[string]string{"A": "1"}, map[string]string{"A": "2"})
	if r.Results[0].Status != "differ" {
		t.Fatalf("expected differ")
	}
}

func TestCompare_LeftOnly(t *testing.T) {
	c := envcompare.New("dev", "staging")
	r := c.Compare(map[string]string{"A": "1"}, map[string]string{})
	if r.Results[0].Status != "left_only" {
		t.Fatalf("expected left_only")
	}
}

func TestCompare_RightOnly(t *testing.T) {
	c := envcompare.New("dev", "staging")
	r := c.Compare(map[string]string{}, map[string]string{"B": "2"})
	if r.Results[0].Status != "right_only" {
		t.Fatalf("expected right_only")
	}
}

func TestSummary_Format(t *testing.T) {
	c := envcompare.New("dev", "staging")
	r := c.Compare(
		map[string]string{"A": "1", "B": "x", "C": "3"},
		map[string]string{"A": "1", "B": "y", "D": "4"},
	)
	s := r.Summary()
	if s == "" {
		t.Fatal("empty summary")
	}
	expected := "match=1 differ=1 left_only=1 right_only=1"
	if s != expected {
		t.Fatalf("got %q want %q", s, expected)
	}
}

func TestLabels_StoredInReport(t *testing.T) {
	c := envcompare.New("alpha", "beta")
	r := c.Compare(nil, nil)
	if r.Left != "alpha" || r.Right != "beta" {
		t.Fatalf("unexpected labels: %s %s", r.Left, r.Right)
	}
}
