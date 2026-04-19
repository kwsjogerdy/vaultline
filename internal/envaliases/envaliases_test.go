package envaliases_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/vaultline/vaultline/internal/envaliases"
)

func TestApply_AddsAliases(t *testing.T) {
	a := envaliases.New(envaliases.AliasMap{
		"DB_PASSWORD": {"DATABASE_PASS", "DB_PASS"},
	})
	secrets := map[string]string{"DB_PASSWORD": "secret123"}
	out := a.Apply(secrets)
	if out["DATABASE_PASS"] != "secret123" {
		t.Errorf("expected DATABASE_PASS=secret123, got %q", out["DATABASE_PASS"])
	}
	if out["DB_PASS"] != "secret123" {
		t.Errorf("expected DB_PASS=secret123, got %q", out["DB_PASS"])
	}
	if out["DB_PASSWORD"] != "secret123" {
		t.Error("original key should be preserved")
	}
}

func TestApply_SkipsMissingSource(t *testing.T) {
	a := envaliases.New(envaliases.AliasMap{
		"MISSING_KEY": {"ALIAS_KEY"},
	})
	out := a.Apply(map[string]string{"OTHER": "val"})
	if _, ok := out["ALIAS_KEY"]; ok {
		t.Error("alias for missing source should not appear")
	}
}

func TestStrip_RemovesAliasedKeys(t *testing.T) {
	a := envaliases.New(envaliases.AliasMap{
		"DB_PASSWORD": {"DATABASE_PASS"},
	})
	secrets := map[string]string{
		"DB_PASSWORD":   "secret",
		"DATABASE_PASS": "secret",
		"API_KEY":       "key",
	}
	out := a.Strip(secrets)
	if _, ok := out["DATABASE_PASS"]; ok {
		t.Error("aliased key should be stripped")
	}
	if out["DB_PASSWORD"] != "secret" {
		t.Error("source key should remain")
	}
	if out["API_KEY"] != "key" {
		t.Error("unrelated key should remain")
	}
}

func TestLoadFile_ParsesJSON(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "aliases.json")
	data, _ := json.Marshal(envaliases.AliasMap{
		"TOKEN": {"AUTH_TOKEN"},
	})
	os.WriteFile(tmp, data, 0600)

	m, err := envaliases.LoadFile(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m["TOKEN"]) != 1 || m["TOKEN"][0] != "AUTH_TOKEN" {
		t.Errorf("unexpected map: %v", m)
	}
}

func TestLoadFile_MissingFile(t *testing.T) {
	_, err := envaliases.LoadFile("/nonexistent/aliases.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
