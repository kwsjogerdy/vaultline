package envfence_test

import (
	"sort"
	"testing"

	"github.com/vaultline/vaultline/internal/envfence"
)

var sample = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PASSWORD": "secret",
	"API_KEY":     "abc123",
	"DEBUG":       "true",
}

func TestApply_AllowPolicy_KeepsOnlyListed(t *testing.T) {
	f, err := envfence.New(envfence.PolicyAllow, []string{"DB_HOST", "DEBUG"})
	if err != nil {
		t.Fatal(err)
	}
	out := f.Apply(sample)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected DB_HOST")
	}
	if _, ok := out["DEBUG"]; !ok {
		t.Error("expected DEBUG")
	}
}

func TestApply_DenyPolicy_RemovesListed(t *testing.T) {
	f, _ := envfence.New(envfence.PolicyDeny, []string{"DB_PASSWORD", "API_KEY"})
	out := f.Apply(sample)
	if _, ok := out["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should be blocked")
	}
	if _, ok := out["API_KEY"]; ok {
		t.Error("API_KEY should be blocked")
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}

func TestApply_EmptyKeyList_AllowDeniesAll(t *testing.T) {
	f, _ := envfence.New(envfence.PolicyAllow, nil)
	out := f.Apply(sample)
	if len(out) != 0 {
		t.Errorf("expected empty result, got %d keys", len(out))
	}
}

func TestBlocked_ReturnsCorrectKeys(t *testing.T) {
	f, _ := envfence.New(envfence.PolicyDeny, []string{"DB_PASSWORD", "API_KEY"})
	blocked := f.Blocked(sample)
	sort.Strings(blocked)
	if len(blocked) != 2 {
		t.Fatalf("expected 2 blocked, got %d", len(blocked))
	}
	if blocked[0] != "API_KEY" || blocked[1] != "DB_PASSWORD" {
		t.Errorf("unexpected blocked keys: %v", blocked)
	}
}

func TestNew_UnknownPolicy_ReturnsError(t *testing.T) {
	_, err := envfence.New("unknown", nil)
	if err == nil {
		t.Error("expected error for unknown policy")
	}
}
