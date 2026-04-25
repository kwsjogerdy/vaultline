package envwildcard_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envwildcard"
)

var sample = map[string]string{
	"DB_HOST":       "localhost",
	"DB_PASSWORD":   "secret",
	"APP_API_KEY":   "key123",
	"APP_DEBUG":     "true",
	"REDIS_URL":     "redis://localhost",
	"INTERNAL_TOKEN": "tok",
}

func TestApply_NoRules_ReturnsAll(t *testing.T) {
	m, err := envwildcard.New(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	out, err := m.Apply(sample)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != len(sample) {
		t.Fatalf("expected %d keys, got %d", len(sample), len(out))
	}
}

func TestApply_IncludePattern(t *testing.T) {
	m, _ := envwildcard.New([]string{"DB_*"}, nil)
	out, err := m.Apply(sample)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 DB_ keys, got %d", len(out))
	}
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected DB_HOST")
	}
}

func TestApply_ExcludePattern(t *testing.T) {
	m, _ := envwildcard.New(nil, []string{"*_PASSWORD", "*_TOKEN"})
	out, err := m.Apply(sample)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := out["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should be excluded")
	}
	if _, ok := out["INTERNAL_TOKEN"]; ok {
		t.Error("INTERNAL_TOKEN should be excluded")
	}
}

func TestApply_IncludeAndExclude(t *testing.T) {
	m, _ := envwildcard.New([]string{"APP_*"}, []string{"*_KEY"})
	out, err := m.Apply(sample)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := out["APP_API_KEY"]; ok {
		t.Error("APP_API_KEY should be excluded")
	}
	if _, ok := out["APP_DEBUG"]; !ok {
		t.Error("APP_DEBUG should be included")
	}
}

func TestNew_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := envwildcard.New([]string{"[invalid"}, nil)
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestMatched_ReturnsSortedKeys(t *testing.T) {
	m, _ := envwildcard.New([]string{"DB_*"}, nil)
	keys, err := m.Matched(sample)
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
	if keys[0] != "DB_HOST" || keys[1] != "DB_PASSWORD" {
		t.Errorf("unexpected order: %v", keys)
	}
}

func TestApply_EmptySecrets(t *testing.T) {
	m, _ := envwildcard.New([]string{"DB_*"}, nil)
	out, err := m.Apply(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map")
	}
}
