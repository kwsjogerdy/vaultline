package envsentinel_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envsentinel"
)

func rules(t *testing.T, rs []envsentinel.Rule) *envsentinel.Sentinel {
	t.Helper()
	s, err := envsentinel.New(rs)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s
}

func TestCheck_NoViolations(t *testing.T) {
	s := rules(t, []envsentinel.Rule{
		{Name: "no-empty-token", Pattern: `TOKEN`, ForbidEmpty: true},
	})
	vs := s.Check(map[string]string{"API_TOKEN": "secret"})
	if len(vs) != 0 {
		t.Fatalf("expected no violations, got %v", vs)
	}
}

func TestCheck_ForbidEmpty(t *testing.T) {
	s := rules(t, []envsentinel.Rule{
		{Name: "no-empty-token", Pattern: `TOKEN`, ForbidEmpty: true},
	})
	vs := s.Check(map[string]string{"API_TOKEN": ""})
	if len(vs) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(vs))
	}
	if vs[0].Key != "API_TOKEN" {
		t.Errorf("expected key API_TOKEN, got %s", vs[0].Key)
	}
}

func TestCheck_ForbidPattern(t *testing.T) {
	s := rules(t, []envsentinel.Rule{
		{Name: "no-placeholder", Pattern: `.*`, ForbidPattern: `^CHANGEME$`},
	})
	vs := s.Check(map[string]string{
		"DB_PASS": "CHANGEME",
		"APP_KEY": "real-value",
	})
	if len(vs) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(vs))
	}
	if vs[0].Key != "DB_PASS" {
		t.Errorf("unexpected key %s", vs[0].Key)
	}
}

func TestCheck_MultipleRules(t *testing.T) {
	s := rules(t, []envsentinel.Rule{
		{Name: "no-empty", Pattern: `SECRET`, ForbidEmpty: true},
		{Name: "no-default", Pattern: `SECRET`, ForbidPattern: `^default$`},
	})
	vs := s.Check(map[string]string{"MY_SECRET": "default"})
	if len(vs) != 1 {
		t.Fatalf("expected 1 violation (forbid pattern), got %d", len(vs))
	}
}

func TestNew_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := envsentinel.New([]envsentinel.Rule{
		{Name: "bad", Pattern: `[invalid`},
	})
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestNew_InvalidForbidPattern_ReturnsError(t *testing.T) {
	_, err := envsentinel.New([]envsentinel.Rule{
		{Name: "bad", Pattern: `.*`, ForbidPattern: `[invalid`},
	})
	if err == nil {
		t.Fatal("expected error for invalid forbid_pattern")
	}
}

func TestHasViolations_True(t *testing.T) {
	s := rules(t, []envsentinel.Rule{
		{Name: "r", Pattern: `KEY`, ForbidEmpty: true},
	})
	if !s.HasViolations(map[string]string{"MY_KEY": ""}) {
		t.Error("expected HasViolations to return true")
	}
}

func TestHasViolations_False(t *testing.T) {
	s := rules(t, []envsentinel.Rule{
		{Name: "r", Pattern: `KEY`, ForbidEmpty: true},
	})
	if s.HasViolations(map[string]string{"MY_KEY": "ok"}) {
		t.Error("expected HasViolations to return false")
	}
}
