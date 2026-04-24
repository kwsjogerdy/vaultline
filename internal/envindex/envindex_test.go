package envindex_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envindex"
)

func TestSearch_ByKey(t *testing.T) {
	idx := envindex.New(map[string]string{
		"DATABASE_URL": "postgres://localhost",
		"API_KEY":      "secret123",
		"REDIS_HOST":   "localhost",
	}, nil)

	results := idx.Search("api")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "API_KEY" {
		t.Errorf("expected API_KEY, got %s", results[0].Key)
	}
	if results[0].MatchedOn != "key" {
		t.Errorf("expected match on key, got %s", results[0].MatchedOn)
	}
}

func TestSearch_ByValue(t *testing.T) {
	idx := envindex.New(map[string]string{
		"DATABASE_URL": "postgres://localhost",
		"REDIS_HOST":   "redis://cache",
	}, nil)

	results := idx.Search("postgres")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].MatchedOn != "value" {
		t.Errorf("expected match on value, got %s", results[0].MatchedOn)
	}
}

func TestSearch_ByLabel(t *testing.T) {
	labels := map[string]string{"SECRET_KEY": "auth"}
	idx := envindex.New(map[string]string{
		"SECRET_KEY": "abc",
		"OTHER":      "xyz",
	}, labels)

	results := idx.Search("auth")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "SECRET_KEY" {
		t.Errorf("unexpected key %s", results[0].Key)
	}
	if results[0].MatchedOn != "label" {
		t.Errorf("expected label match, got %s", results[0].MatchedOn)
	}
}

func TestSearch_NoMatch(t *testing.T) {
	idx := envindex.New(map[string]string{"FOO": "bar"}, nil)
	results := idx.Search("zzz")
	if len(results) != 0 {
		t.Errorf("expected no results, got %d", len(results))
	}
}

func TestSearch_CaseInsensitive(t *testing.T) {
	idx := envindex.New(map[string]string{"MY_TOKEN": "abc"}, nil)
	results := idx.Search("TOKEN")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestKeys_SortedOrder(t *testing.T) {
	idx := envindex.New(map[string]string{
		"ZEBRA": "1",
		"ALPHA": "2",
		"MANGO": "3",
	}, nil)
	keys := idx.Keys()
	expected := []string{"ALPHA", "MANGO", "ZEBRA"}
	for i, k := range expected {
		if keys[i] != k {
			t.Errorf("position %d: expected %s, got %s", i, k, keys[i])
		}
	}
}

func TestSize(t *testing.T) {
	idx := envindex.New(map[string]string{"A": "1", "B": "2"}, nil)
	if idx.Size() != 2 {
		t.Errorf("expected size 2, got %d", idx.Size())
	}
}
