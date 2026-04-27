package envreplace_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envreplace"
)

func TestNew_EmptyFind_ReturnsError(t *testing.T) {
	_, err := envreplace.New([]envreplace.Rule{{Find: "", Replace: "x"}})
	if err == nil {
		t.Fatal("expected error for empty Find")
	}
}

func TestApply_AllKeys(t *testing.T) {
	r, err := envreplace.New([]envreplace.Rule{{Find: "localhost", Replace: "prod.example.com"}})
	if err != nil {
		t.Fatal(err)
	}
	secrets := map[string]string{
		"DB_HOST": "localhost:5432",
		"REDIS":   "redis://localhost",
		"APP":     "myapp",
	}
	out := r.Apply(secrets)
	if out["DB_HOST"] != "prod.example.com:5432" {
		t.Errorf("DB_HOST: got %q", out["DB_HOST"])
	}
	if out["REDIS"] != "redis://prod.example.com" {
		t.Errorf("REDIS: got %q", out["REDIS"])
	}
	if out["APP"] != "myapp" {
		t.Errorf("APP should be unchanged, got %q", out["APP"])
	}
}

func TestApply_TargetedKeys(t *testing.T) {
	r, _ := envreplace.New([]envreplace.Rule{
		{Find: "dev", Replace: "prod", Keys: []string{"ENV_NAME"}},
	})
	secrets := map[string]string{
		"ENV_NAME":  "dev-cluster",
		"OTHER_ENV": "dev-other",
	}
	out := r.Apply(secrets)
	if out["ENV_NAME"] != "prod-cluster" {
		t.Errorf("ENV_NAME: got %q", out["ENV_NAME"])
	}
	if out["OTHER_ENV"] != "dev-other" {
		t.Errorf("OTHER_ENV should be unchanged, got %q", out["OTHER_ENV"])
	}
}

func TestApply_MultipleRules(t *testing.T) {
	r, _ := envreplace.New([]envreplace.Rule{
		{Find: "foo", Replace: "bar"},
		{Find: "bar", Replace: "baz"},
	})
	secrets := map[string]string{"KEY": "foo"}
	out := r.Apply(secrets)
	if out["KEY"] != "baz" {
		t.Errorf("expected chained replacement, got %q", out["KEY"])
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	r, _ := envreplace.New([]envreplace.Rule{{Find: "old", Replace: "new"}})
	secrets := map[string]string{"K": "old-value"}
	r.Apply(secrets)
	if secrets["K"] != "old-value" {
		t.Error("input map was mutated")
	}
}

func TestChanged_ReturnsModifiedKeys(t *testing.T) {
	r, _ := envreplace.New([]envreplace.Rule{{Find: "http", Replace: "https"}})
	secrets := map[string]string{
		"URL":  "http://example.com",
		"NAME": "myapp",
	}
	changed := r.Changed(secrets)
	if len(changed) != 1 || changed[0] != "URL" {
		t.Errorf("expected [URL], got %v", changed)
	}
}

func TestApply_EmptySecrets(t *testing.T) {
	r, _ := envreplace.New([]envreplace.Rule{{Find: "x", Replace: "y"}})
	out := r.Apply(map[string]string{})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}
