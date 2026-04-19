package envsort

import (
	"testing"
)

func TestSort_ByKeyAsc(t *testing.T) {
	s := New(Options{By: ByKey, Order: Asc})
	secrets := map[string]string{"ZEBRA": "1", "ALPHA": "2", "MANGO": "3"}
	pairs := s.Sort(secrets)
	if pairs[0][0] != "ALPHA" || pairs[1][0] != "MANGO" || pairs[2][0] != "ZEBRA" {
		t.Fatalf("unexpected order: %v", pairs)
	}
}

func TestSort_ByKeyDesc(t *testing.T) {
	s := New(Options{By: ByKey, Order: Desc})
	secrets := map[string]string{"ZEBRA": "1", "ALPHA": "2", "MANGO": "3"}
	pairs := s.Sort(secrets)
	if pairs[0][0] != "ZEBRA" || pairs[2][0] != "ALPHA" {
		t.Fatalf("unexpected order: %v", pairs)
	}
}

func TestSort_ByValueAsc(t *testing.T) {
	s := New(Options{By: ByValue, Order: Asc})
	secrets := map[string]string{"A": "charlie", "B": "alpha", "C": "bravo"}
	pairs := s.Sort(secrets)
	if pairs[0][1] != "alpha" || pairs[1][1] != "bravo" || pairs[2][1] != "charlie" {
		t.Fatalf("unexpected order: %v", pairs)
	}
}

func TestSort_Defaults(t *testing.T) {
	s := New(Options{})
	if s.opts.By != ByKey {
		t.Errorf("expected default By=key, got %s", s.opts.By)
	}
	if s.opts.Order != Asc {
		t.Errorf("expected default Order=asc, got %s", s.opts.Order)
	}
}

func TestKeys_ReturnsSortedKeys(t *testing.T) {
	s := New(Options{By: ByKey, Order: Asc})
	secrets := map[string]string{"Z": "z", "A": "a", "M": "m"}
	keys := s.Keys(secrets)
	if keys[0] != "A" || keys[1] != "M" || keys[2] != "Z" {
		t.Fatalf("unexpected keys: %v", keys)
	}
}

func TestSort_EmptyMap(t *testing.T) {
	s := New(Options{})
	pairs := s.Sort(map[string]string{})
	if len(pairs) != 0 {
		t.Errorf("expected empty result")
	}
}
