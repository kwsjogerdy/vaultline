package envpatch_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envpatch"
)

func base() map[string]string {
	return map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"API_KEY": "secret",
	}
}

func TestApply_SetAddsAndOverwrites(t *testing.T) {
	p, _ := envpatch.New([]envpatch.Op{
		{Type: envpatch.OpSet, Key: "DB_HOST", Value: "prod.db"},
		{Type: envpatch.OpSet, Key: "NEW_KEY", Value: "hello"},
	})
	out, res, err := p.Apply(base())
	if err != nil {
		t.Fatal(err)
	}
	if out["DB_HOST"] != "prod.db" {
		t.Errorf("expected prod.db, got %s", out["DB_HOST"])
	}
	if out["NEW_KEY"] != "hello" {
		t.Errorf("expected hello, got %s", out["NEW_KEY"])
	}
	if len(res.Applied) != 2 {
		t.Errorf("expected 2 applied, got %d", len(res.Applied))
	}
}

func TestApply_DeleteRemovesKey(t *testing.T) {
	p, _ := envpatch.New([]envpatch.Op{
		{Type: envpatch.OpDelete, Key: "API_KEY"},
	})
	out, _, err := p.Apply(base())
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := out["API_KEY"]; ok {
		t.Error("expected API_KEY to be deleted")
	}
}

func TestApply_DeleteMissingKeySkips(t *testing.T) {
	p, _ := envpatch.New([]envpatch.Op{
		{Type: envpatch.OpDelete, Key: "GHOST"},
	})
	_, res, err := p.Apply(base())
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "GHOST" {
		t.Errorf("expected GHOST in skipped, got %v", res.Skipped)
	}
}

func TestApply_RenameMovesKey(t *testing.T) {
	p, _ := envpatch.New([]envpatch.Op{
		{Type: envpatch.OpRename, Key: "DB_PORT", NewKey: "DATABASE_PORT"},
	})
	out, _, err := p.Apply(base())
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := out["DB_PORT"]; ok {
		t.Error("old key should be removed")
	}
	if out["DATABASE_PORT"] != "5432" {
		t.Errorf("expected 5432, got %s", out["DATABASE_PORT"])
	}
}

func TestNew_EmptyKeyReturnsError(t *testing.T) {
	_, err := envpatch.New([]envpatch.Op{
		{Type: envpatch.OpSet, Key: "", Value: "x"},
	})
	if err == nil {
		t.Error("expected error for empty key")
	}
}

func TestNew_RenameWithoutNewKeyReturnsError(t *testing.T) {
	_, err := envpatch.New([]envpatch.Op{
		{Type: envpatch.OpRename, Key: "FOO", NewKey: ""},
	})
	if err == nil {
		t.Error("expected error for missing new_key")
	}
}

func TestNew_UnknownOpReturnsError(t *testing.T) {
	_, err := envpatch.New([]envpatch.Op{
		{Type: "upsert", Key: "FOO"},
	})
	if err == nil {
		t.Error("expected error for unknown op type")
	}
}
