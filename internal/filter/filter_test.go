package filter_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/filter"
)

var sampleSecrets = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PASSWORD": "secret",
	"APP_KEY":     "abc123",
	"INTERNAL_ID": "skip-me",
}

func TestApply_NoRules_ReturnsAll(t *testing.T) {
	f := filter.New(filter.Rule{})
	result := f.Apply(sampleSecrets)
	if len(result) != len(sampleSecrets) {
		t.Errorf("expected %d keys, got %d", len(sampleSecrets), len(result))
	}
}

func TestApply_PrefixFilter(t *testing.T) {
	f := filter.New(filter.Rule{Prefixes: []string{"DB_"}})
	result := f.Apply(sampleSecrets)
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
	if _, ok := result["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in result")
	}
}

func TestApply_ExcludeFilter(t *testing.T) {
	f := filter.New(filter.Rule{Excludes: []string{"INTERNAL_ID"}})
	result := f.Apply(sampleSecrets)
	if _, ok := result["INTERNAL_ID"]; ok {
		t.Error("expected INTERNAL_ID to be excluded")
	}
	if len(result) != 3 {
		t.Errorf("expected 3 keys, got %d", len(result))
	}
}

func TestApply_PrefixAndExclude(t *testing.T) {
	f := filter.New(filter.Rule{
		Prefixes: []string{"DB_"},
		Excludes: []string{"DB_PASSWORD"},
	})
	result := f.Apply(sampleSecrets)
	if len(result) != 1 {
		t.Errorf("expected 1 key, got %d", len(result))
	}
	if _, ok := result["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in result")
	}
}

func TestApply_EmptySecrets(t *testing.T) {
	f := filter.New(filter.Rule{Prefixes: []string{"DB_"}})
	result := f.Apply(map[string]string{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d keys", len(result))
	}
}

func TestApply_ExcludeNonExistentKey(t *testing.T) {
	f := filter.New(filter.Rule{Excludes: []string{"DOES_NOT_EXIST"}})
	result := f.Apply(sampleSecrets)
	if len(result) != len(sampleSecrets) {
		t.Errorf("expected %d keys, got %d", len(sampleSecrets), len(result))
	}
}
