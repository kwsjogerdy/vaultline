package envprefix_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envprefix"
)

func TestNew_BothEmpty_ReturnsError(t *testing.T) {
	_, err := envprefix.New("", "")
	if err == nil {
		t.Fatal("expected error for empty prefixes")
	}
}

func TestApply_AddPrefix(t *testing.T) {
	tr, _ := envprefix.New("APP_", "")
	result := tr.Apply(map[string]string{"DB_HOST": "localhost", "PORT": "5432"})
	if result["APP_DB_HOST"] != "localhost" {
		t.Errorf("expected APP_DB_HOST, got %v", result)
	}
	if result["APP_PORT"] != "5432" {
		t.Errorf("expected APP_PORT, got %v", result)
	}
}

func TestApply_RemovePrefix(t *testing.T) {
	tr, _ := envprefix.New("", "APP_")
	result := tr.Apply(map[string]string{"APP_DB_HOST": "localhost", "PORT": "5432"})
	if result["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST, got %v", result)
	}
	if result["PORT"] != "5432" {
		t.Errorf("expected PORT, got %v", result)
	}
}

func TestApply_RemoveThenAdd(t *testing.T) {
	tr, _ := envprefix.New("NEW_", "OLD_")
	result := tr.Apply(map[string]string{"OLD_KEY": "val"})
	if result["NEW_KEY"] != "val" {
		t.Errorf("expected NEW_KEY=val, got %v", result)
	}
}

func TestApply_SkipsEmptyResultKey(t *testing.T) {
	tr, _ := envprefix.New("", "APP_")
	result := tr.Apply(map[string]string{"APP_": "orphan"})
	if _, ok := result[""]; ok {
		t.Error("empty key should be dropped")
	}
}

func TestReplacePrefix(t *testing.T) {
	secrets := map[string]string{
		"DEV_DB":  "devdb",
		"DEV_API": "devapi",
		"SHARED":  "common",
	}
	result := envprefix.ReplacePrefix(secrets, "DEV_", "PROD_")
	if result["PROD_DB"] != "devdb" {
		t.Errorf("expected PROD_DB, got %v", result)
	}
	if result["PROD_API"] != "devapi" {
		t.Errorf("expected PROD_API, got %v", result)
	}
	if result["SHARED"] != "common" {
		t.Errorf("expected SHARED unchanged, got %v", result)
	}
}

func TestReplacePrefix_NoMatch(t *testing.T) {
	secrets := map[string]string{"KEY": "val"}
	result := envprefix.ReplacePrefix(secrets, "MISSING_", "NEW_")
	if result["KEY"] != "val" {
		t.Errorf("expected KEY unchanged, got %v", result)
	}
}
