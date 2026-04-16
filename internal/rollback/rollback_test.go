package rollback

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListBackups_ReturnsNewestFirst(t *testing.T) {
	dir := t.TempDir()
	names := []string{".env.20240101_120000", ".env.20240102_120000", ".env.20240103_120000"}
	for _, n := range names {
		if err := os.WriteFile(filepath.Join(dir, n), []byte("x"), 0600); err != nil {
			t.Fatal(err)
		}
	}

	m := New(dir)
	results, err := m.ListBackups(".env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 backups, got %d", len(results))
	}
	if !contains(results[0], "20240103") {
		t.Errorf("expected newest first, got %s", results[0])
	}
}

func TestListBackups_MissingDir(t *testing.T) {
	m := New("/nonexistent/dir")
	results, err := m.ListBackups(".env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results != nil {
		t.Errorf("expected nil, got %v", results)
	}
}

func TestRollback_RestoresMostRecent(t *testing.T) {
	dir := t.TempDir()
	backup := filepath.Join(dir, ".env.20240103_120000")
	if err := os.WriteFile(backup, []byte("SECRET=restored\n"), 0600); err != nil {
		t.Fatal(err)
	}

	envFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(envFile, []byte("SECRET=old\n"), 0600); err != nil {
		t.Fatal(err)
	}

	m := New(dir)
	used, err := m.Rollback(envFile, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if used != backup {
		t.Errorf("expected %s, got %s", backup, used)
	}

	data, _ := os.ReadFile(envFile)
	if string(data) != "SECRET=restored\n" {
		t.Errorf("unexpected content: %s", data)
	}
}

func TestRollback_NoBackups(t *testing.T) {
	dir := t.TempDir()
	m := New(dir)
	_, err := m.Rollback(filepath.Join(dir, ".env"), "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
