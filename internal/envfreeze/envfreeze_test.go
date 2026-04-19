package envfreeze_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vaultline/vaultline/internal/envfreeze"
)

func tmpFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "freeze.json")
}

func TestFreeze_AndIsFrozen(t *testing.T) {
	f, err := envfreeze.New(tmpFile(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Freeze("DB_PASSWORD"); err != nil {
		t.Fatal(err)
	}
	if !f.IsFrozen("DB_PASSWORD") {
		t.Error("expected DB_PASSWORD to be frozen")
	}
	if f.IsFrozen("API_KEY") {
		t.Error("API_KEY should not be frozen")
	}
}

func TestUnfreeze_RemovesKey(t *testing.T) {
	f, _ := envfreeze.New(tmpFile(t))
	_ = f.Freeze("TOKEN")
	_ = f.Unfreeze("TOKEN")
	if f.IsFrozen("TOKEN") {
		t.Error("TOKEN should no longer be frozen")
	}
}

func TestFreeze_EmptyKey_ReturnsError(t *testing.T) {
	f, _ := envfreeze.New(tmpFile(t))
	if err := f.Freeze(""); err == nil {
		t.Error("expected error for empty key")
	}
}

func TestList_SortedByKey(t *testing.T) {
	f, _ := envfreeze.New(tmpFile(t))
	_ = f.Freeze("Z_KEY")
	_ = f.Freeze("A_KEY")
	list := f.List()
	if len(list) != 2 || list[0].Key != "A_KEY" || list[1].Key != "Z_KEY" {
		t.Errorf("unexpected order: %v", list)
	}
}

func TestFilter_RemovesFrozenKeys(t *testing.T) {
	f, _ := envfreeze.New(tmpFile(t))
	_ = f.Freeze("SECRET")
	secrets := map[string]string{"SECRET": "s3cr3t", "SAFE": "ok"}
	result := f.Filter(secrets)
	if _, ok := result["SECRET"]; ok {
		t.Error("frozen key should be excluded")
	}
	if result["SAFE"] != "ok" {
		t.Error("non-frozen key should be retained")
	}
}

func TestPersistence_ReloadsFromDisk(t *testing.T) {
	path := tmpFile(t)
	f, _ := envfreeze.New(path)
	_ = f.Freeze("PERSIST_KEY")

	f2, err := envfreeze.New(path)
	if err != nil {
		t.Fatal(err)
	}
	if !f2.IsFrozen("PERSIST_KEY") {
		t.Error("expected PERSIST_KEY to survive reload")
	}
}

func TestNew_MissingFile_Empty(t *testing.T) {
	f, err := envfreeze.New(filepath.Join(t.TempDir(), "nonexistent.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(f.List()) != 0 {
		t.Error("expected empty list for missing file")
	}
}

func TestNew_CorruptFile_ReturnsError(t *testing.T) {
	path := tmpFile(t)
	_ = os.WriteFile(path, []byte("not json"), 0o600)
	if _, err := envfreeze.New(path); err == nil {
		t.Error("expected error for corrupt file")
	}
}
