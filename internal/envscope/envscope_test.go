package envscope_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envscope"
)

var secrets = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PASSWORD": "secret",
	"APP_KEY":     "abc123",
	"CACHE_URL":   "redis://localhost",
}

func TestApply_UnknownScope(t *testing.T) {
	m := envscope.New(nil)
	_, err := m.Apply("prod", secrets)
	if err == nil {
		t.Fatal("expected error for unknown scope")
	}
}

func TestApply_NoPrefixes_ReturnsAll(t *testing.T) {
	m := envscope.New([]envscope.Scope{{Name: "dev", Prefixes: nil}})
	got, err := m.Apply("dev", secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(secrets) {
		t.Fatalf("expected %d keys, got %d", len(secrets), len(got))
	}
}

func TestApply_PrefixFilter(t *testing.T) {
	m := envscope.New([]envscope.Scope{
		{Name: "db-only", Prefixes: []string{"DB_"}},
	})
	got, err := m.Apply("db-only", secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(got))
	}
	if _, ok := got["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in result")
	}
}

func TestApply_MultiplePrefix(t *testing.T) {
	m := envscope.New([]envscope.Scope{
		{Name: "infra", Prefixes: []string{"DB_", "CACHE_"}},
	})
	got, err := m.Apply("infra", secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(got))
	}
}

func TestApply_CaseInsensitiveScopeName(t *testing.T) {
	m := envscope.New([]envscope.Scope{{Name: "Staging", Prefixes: []string{"APP_"}}})
	got, err := m.Apply("STAGING", secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 key, got %d", len(got))
	}
}

func TestList_ReturnsScopeNames(t *testing.T) {
	m := envscope.New([]envscope.Scope{
		{Name: "dev"},
		{Name: "prod"},
	})
	names := m.List()
	if len(names) != 2 {
		t.Fatalf("expected 2 scopes, got %d", len(names))
	}
}
