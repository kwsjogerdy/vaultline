package envlabel_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vaultline/vaultline/internal/envlabel"
)

func tmpFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "labels.json")
}

func TestNew_EmptyWhenMissing(t *testing.T) {
	l, err := envlabel.New(tmpFile(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if keys := l.Filter("env", "prod"); len(keys) != 0 {
		t.Errorf("expected no keys, got %v", keys)
	}
}

func TestSet_AndGet(t *testing.T) {
	l, _ := envlabel.New(tmpFile(t))
	l.Set("DB_PASSWORD", "env", "prod")
	v, ok := l.Get("DB_PASSWORD", "env")
	if !ok {
		t.Fatal("expected label to exist")
	}
	if v != "prod" {
		t.Errorf("got %q, want %q", v, "prod")
	}
}

func TestGet_MissingKey(t *testing.T) {
	l, _ := envlabel.New(tmpFile(t))
	_, ok := l.Get("NONEXISTENT", "env")
	if ok {
		t.Error("expected false for missing key")
	}
}

func TestFilter_ReturnsSortedMatches(t *testing.T) {
	l, _ := envlabel.New(tmpFile(t))
	l.Set("API_KEY", "env", "prod")
	l.Set("DB_PASS", "env", "prod")
	l.Set("DEV_TOKEN", "env", "dev")

	got := l.Filter("env", "prod")
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
	if got[0] != "API_KEY" || got[1] != "DB_PASS" {
		t.Errorf("unexpected order: %v", got)
	}
}

func TestFilter_NoMatches(t *testing.T) {
	l, _ := envlabel.New(tmpFile(t))
	l.Set("API_KEY", "env", "dev")
	got := l.Filter("env", "prod")
	if len(got) != 0 {
		t.Errorf("expected empty, got %v", got)
	}
}

func TestSaveAndReload(t *testing.T) {
	path := tmpFile(t)
	l, _ := envlabel.New(path)
	l.Set("SECRET_A", "team", "backend")
	l.Set("SECRET_B", "team", "frontend")
	l.Set("SECRET_A", "env", "staging")
	if err := l.Save(); err != nil {
		t.Fatalf("save: %v", err)
	}

	l2, err := envlabel.New(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	v, ok := l2.Get("SECRET_A", "team")
	if !ok || v != "backend" {
		t.Errorf("got %q %v, want backend true", v, ok)
	}
	v2, ok2 := l2.Get("SECRET_A", "env")
	if !ok2 || v2 != "staging" {
		t.Errorf("got %q %v, want staging true", v2, ok2)
	}
	keys := l2.Filter("team", "frontend")
	if len(keys) != 1 || keys[0] != "SECRET_B" {
		t.Errorf("unexpected filter result: %v", keys)
	}
}

func TestNew_InvalidJSON(t *testing.T) {
	path := tmpFile(t)
	_ = os.WriteFile(path, []byte("not-json"), 0o600)
	_, err := envlabel.New(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
