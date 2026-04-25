package envcount_test

import (
	"strings"
	"testing"

	"github.com/vaultline/vaultline/internal/envcount"
)

func TestCompute_EmptyMap(t *testing.T) {
	c := envcount.New()
	s := c.Compute(map[string]string{})
	if s.Total != 0 || s.Empty != 0 || s.NonEmpty != 0 {
		t.Fatalf("expected all zeros, got %+v", s)
	}
}

func TestCompute_CountsTotals(t *testing.T) {
	c := envcount.New()
	secrets := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_PASS":  "",
		"PLAIN":    "value",
	}
	s := c.Compute(secrets)
	if s.Total != 4 {
		t.Errorf("Total: want 4, got %d", s.Total)
	}
	if s.Empty != 1 {
		t.Errorf("Empty: want 1, got %d", s.Empty)
	}
	if s.NonEmpty != 3 {
		t.Errorf("NonEmpty: want 3, got %d", s.NonEmpty)
	}
}

func TestCompute_ByPrefix(t *testing.T) {
	c := envcount.New()
	secrets := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_PASS":  "secret",
		"NOPREFIX": "x",
	}
	s := c.Compute(secrets)
	if s.ByPrefix["APP"] != 2 {
		t.Errorf("APP prefix: want 2, got %d", s.ByPrefix["APP"])
	}
	if s.ByPrefix["DB"] != 1 {
		t.Errorf("DB prefix: want 1, got %d", s.ByPrefix["DB"])
	}
	if _, ok := s.ByPrefix["NOPREFIX"]; ok {
		t.Error("NOPREFIX should not appear in ByPrefix")
	}
}

func TestCompute_WhitespaceValueCountsAsEmpty(t *testing.T) {
	c := envcount.New()
	s := c.Compute(map[string]string{"KEY": "   "})
	if s.Empty != 1 {
		t.Errorf("want 1 empty, got %d", s.Empty)
	}
}

func TestSummary_ContainsExpectedLines(t *testing.T) {
	c := envcount.New()
	secrets := map[string]string{
		"APP_A": "1",
		"APP_B": "",
		"DB_X":  "y",
	}
	s := c.Compute(secrets)
	summary := s.Summary()
	for _, want := range []string{"total:", "empty:", "non-empty:", "by prefix:", "APP:", "DB:"} {
		if !strings.Contains(summary, want) {
			t.Errorf("summary missing %q\n%s", want, summary)
		}
	}
}
