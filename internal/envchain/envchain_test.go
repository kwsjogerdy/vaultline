package envchain_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envchain"
)

func TestResolve_HigherPriorityWins(t *testing.T) {
	chain := envchain.New([]envchain.Source{
		{Name: "local", Priority: 1, Secrets: map[string]string{"DB_HOST": "localhost", "APP_ENV": "dev"}},
		{Name: "vault", Priority: 10, Secrets: map[string]string{"DB_HOST": "vault-host", "DB_PASS": "secret"}},
	})
	resolved := chain.Resolve()
	if resolved["DB_HOST"] != "vault-host" {
		t.Errorf("expected vault-host, got %s", resolved["DB_HOST"])
	}
	if resolved["APP_ENV"] != "dev" {
		t.Errorf("expected dev, got %s", resolved["APP_ENV"])
	}
	if resolved["DB_PASS"] != "secret" {
		t.Errorf("expected secret, got %s", resolved["DB_PASS"])
	}
}

func TestResolve_SingleSource(t *testing.T) {
	chain := envchain.New([]envchain.Source{
		{Name: "only", Priority: 5, Secrets: map[string]string{"KEY": "val"}},
	})
	if chain.Resolve()["KEY"] != "val" {
		t.Error("expected val")
	}
}

func TestOrigin_FindsSource(t *testing.T) {
	chain := envchain.New([]envchain.Source{
		{Name: "local", Priority: 1, Secrets: map[string]string{"ONLY_LOCAL": "x"}},
		{Name: "vault", Priority: 10, Secrets: map[string]string{"ONLY_VAULT": "y"}},
	})
	name, err := chain.Origin("ONLY_VAULT")
	if err != nil || name != "vault" {
		t.Errorf("expected vault, got %s %v", name, err)
	}
}

func TestOrigin_NotFound(t *testing.T) {
	chain := envchain.New([]envchain.Source{})
	_, err := chain.Origin("MISSING")
	if err == nil {
		t.Error("expected error for missing key")
	}
}

func TestResolve_EmptySources(t *testing.T) {
	chain := envchain.New([]envchain.Source{})
	if len(chain.Resolve()) != 0 {
		t.Error("expected empty map")
	}
}
