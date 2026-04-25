package envstats_test

import (
	"strings"
	"testing"

	"github.com/vaultline/vaultline/internal/envstats"
)

func TestCompute_EmptyMap(t *testing.T) {
	a := envstats.New("_")
	s := a.Compute(map[string]string{})
	if s.Total != 0 {
		t.Fatalf("expected 0 total, got %d", s.Total)
	}
	if s.Empty != 0 {
		t.Fatalf("expected 0 empty, got %d", s.Empty)
	}
}

func TestCompute_BasicStats(t *testing.T) {
	a := envstats.New("_")
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_KEY": "secret",
		"EMPTY":   "",
	}
	s := a.Compute(secrets)
	if s.Total != 4 {
		t.Fatalf("expected 4 total, got %d", s.Total)
	}
	if s.Empty != 1 {
		t.Fatalf("expected 1 empty, got %d", s.Empty)
	}
	if s.MaxValueLen != 9 { // "localhost"
		t.Fatalf("expected max 9, got %d", s.MaxValueLen)
	}
	if s.MinValueLen != 0 {
		t.Fatalf("expected min 0, got %d", s.MinValueLen)
	}
}

func TestCompute_PrefixCounts(t *testing.T) {
	a := envstats.New("_")
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_KEY": "secret",
	}
	s := a.Compute(secrets)
	if s.PrefixCounts["DB"] != 2 {
		t.Fatalf("expected DB prefix count 2, got %d", s.PrefixCounts["DB"])
	}
	if s.PrefixCounts["APP"] != 1 {
		t.Fatalf("expected APP prefix count 1, got %d", s.PrefixCounts["APP"])
	}
}

func TestCompute_UniqueValues(t *testing.T) {
	a := envstats.New("_")
	secrets := map[string]string{
		"A": "same",
		"B": "same",
		"C": "different",
	}
	s := a.Compute(secrets)
	if s.Unique != 2 {
		t.Fatalf("expected 2 unique values, got %d", s.Unique)
	}
}

func TestCompute_NoDelimiterKey(t *testing.T) {
	a := envstats.New("_")
	secrets := map[string]string{
		"NODASH": "value",
	}
	s := a.Compute(secrets)
	if s.PrefixCounts["(none)"] != 1 {
		t.Fatalf("expected (none) prefix count 1, got %d", s.PrefixCounts["(none)"])
	}
}

func TestSummary_ContainsFields(t *testing.T) {
	a := envstats.New("_")
	s := a.Compute(map[string]string{"APP_KEY": "val"})
	out := s.Summary()
	for _, want := range []string{"Total", "Empty", "Unique", "Avg", "Max", "Min", "Prefix"} {
		if !strings.Contains(out, want) {
			t.Errorf("summary missing field %q", want)
		}
	}
}
