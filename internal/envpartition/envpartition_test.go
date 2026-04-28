package envpartition_test

import (
	"strings"
	"testing"

	"github.com/vaultline/vaultline/internal/envpartition"
)

func TestApply_MatchedAndUnmatched(t *testing.T) {
	p, err := envpartition.New(func(k, _ string) bool {
		return strings.HasPrefix(k, "DB_")
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_ENV":  "production",
	}
	res := p.Apply(secrets)
	if len(res.Matched) != 2 {
		t.Errorf("expected 2 matched, got %d", len(res.Matched))
	}
	if len(res.Unmatched) != 1 {
		t.Errorf("expected 1 unmatched, got %d", len(res.Unmatched))
	}
}

func TestApply_EmptySecrets(t *testing.T) {
	p, _ := envpartition.New(func(k, v string) bool { return true })
	res := p.Apply(map[string]string{})
	if len(res.Matched) != 0 || len(res.Unmatched) != 0 {
		t.Error("expected empty result for empty input")
	}
}

func TestApply_AllMatch(t *testing.T) {
	p, _ := envpartition.New(func(k, v string) bool { return true })
	secrets := map[string]string{"A": "1", "B": "2"}
	res := p.Apply(secrets)
	if len(res.Matched) != 2 || len(res.Unmatched) != 0 {
		t.Errorf("expected all matched, got matched=%d unmatched=%d", len(res.Matched), len(res.Unmatched))
	}
}

func TestApply_NoneMatch(t *testing.T) {
	p, _ := envpartition.New(func(k, v string) bool { return false })
	secrets := map[string]string{"X": "foo", "Y": "bar"}
	res := p.Apply(secrets)
	if len(res.Matched) != 0 || len(res.Unmatched) != 2 {
		t.Errorf("expected none matched, got matched=%d unmatched=%d", len(res.Matched), len(res.Unmatched))
	}
}

func TestNew_NilPredicate_ReturnsError(t *testing.T) {
	_, err := envpartition.New(nil)
	if err == nil {
		t.Fatal("expected error for nil predicate")
	}
}

func TestMatchedKeys_Sorted(t *testing.T) {
	p, _ := envpartition.New(func(k, _ string) bool { return true })
	secrets := map[string]string{"Z": "1", "A": "2", "M": "3"}
	res := p.Apply(secrets)
	keys := res.MatchedKeys()
	if len(keys) != 3 || keys[0] != "A" || keys[1] != "M" || keys[2] != "Z" {
		t.Errorf("expected sorted keys [A M Z], got %v", keys)
	}
}

func TestPredicateReceivesValue(t *testing.T) {
	p, _ := envpartition.New(func(_, v string) bool {
		return strings.Contains(v, "secret")
	})
	secrets := map[string]string{
		"TOKEN": "my-secret-token",
		"ENV":   "staging",
	}
	res := p.Apply(secrets)
	if _, ok := res.Matched["TOKEN"]; !ok {
		t.Error("expected TOKEN in matched bucket")
	}
	if _, ok := res.Unmatched["ENV"]; !ok {
		t.Error("expected ENV in unmatched bucket")
	}
}
