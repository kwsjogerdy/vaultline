package envpin_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vaultline/vaultline/internal/envpin"
)

func tmpFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "vaultline.lock.json")
}

func TestLoad_EmptyWhenMissing(t *testing.T) {
	m := envpin.New(tmpFile(t))
	lf, err := m.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lf.Pins) != 0 {
		t.Errorf("expected empty pins, got %d", len(lf.Pins))
	}
}

func TestSetAndGet(t *testing.T) {
	m := envpin.New(tmpFile(t))
	lf, _ := m.Load()
	m.Set(lf, "secret/app", 3)
	v, ok := m.Get(lf, "secret/app")
	if !ok || v != 3 {
		t.Errorf("expected version 3, got %d ok=%v", v, ok)
	}
}

func TestSaveAndReload(t *testing.T) {
	path := tmpFile(t)
	m := envpin.New(path)
	lf, _ := m.Load()
	m.Set(lf, "secret/db", 7)
	if err := m.Save(lf); err != nil {
		t.Fatalf("save: %v", err)
	}
	m2 := envpin.New(path)
	lf2, err := m2.Load()
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	v, ok := m2.Get(lf2, "secret/db")
	if !ok || v != 7 {
		t.Errorf("expected version 7, got %d ok=%v", v, ok)
	}
}

func TestRemove(t *testing.T) {
	m := envpin.New(tmpFile(t))
	lf, _ := m.Load()
	m.Set(lf, "secret/x", 1)
	m.Remove(lf, "secret/x")
	_, ok := m.Get(lf, "secret/x")
	if ok {
		t.Error("expected pin to be removed")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	path := tmpFile(t)
	_ = os.WriteFile(path, []byte("not-json"), 0600)
	m := envpin.New(path)
	_, err := m.Load()
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
