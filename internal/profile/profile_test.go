package profile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vaultline/vaultline/internal/profile"
)

func tmpStore(t *testing.T) *profile.Store {
	t.Helper()
	dir := t.TempDir()
	return profile.New(filepath.Join(dir, "profiles"))
}

func TestSaveAndLoad(t *testing.T) {
	s := tmpStore(t)
	p := profile.Profile{Name: "dev", VaultPath: "secret/dev", OutputFile: ".env"}
	if err := s.Save(p); err != nil {
		t.Fatal(err)
	}
	got, err := s.Load("dev")
	if err != nil {
		t.Fatal(err)
	}
	if got.VaultPath != p.VaultPath {
		t.Errorf("want %q got %q", p.VaultPath, got.VaultPath)
	}
}

func TestLoad_NotFound(t *testing.T) {
	s := tmpStore(t)
	_, err := s.Load("missing")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSave_EmptyName(t *testing.T) {
	s := tmpStore(t)
	err := s.Save(profile.Profile{})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestList(t *testing.T) {
	s := tmpStore(t)
	for _, name := range []string{"dev", "staging", "prod"} {
		_ = s.Save(profile.Profile{Name: name, VaultPath: "secret/" + name})
	}
	names, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 3 {
		t.Errorf("want 3 profiles, got %d", len(names))
	}
}

func TestList_EmptyDir(t *testing.T) {
	s := tmpStore(t)
	names, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 0 {
		t.Errorf("expected empty list")
	}
}

func TestDelete(t *testing.T) {
	s := tmpStore(t)
	_ = s.Save(profile.Profile{Name: "dev", VaultPath: "secret/dev"})
	if err := s.Delete("dev"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(t.TempDir(), "profiles", "dev.json")); !os.IsNotExist(err) {
		t.Log("file removed as expected")
	}
}

func TestDelete_NotFound(t *testing.T) {
	s := tmpStore(t)
	if err := s.Delete("ghost"); err == nil {
		t.Fatal("expected error")
	}
}
