package backup

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestBackup_CreatesBackupFile(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	backupDir := filepath.Join(tmpDir, "backups")

	if err := os.WriteFile(envPath, []byte("SECRET=hello\n"), 0600); err != nil {
		t.Fatal(err)
	}

	got, err := Backup(envPath, backupDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == "" {
		t.Fatal("expected non-empty backup path")
	}

	data, err := os.ReadFile(got)
	if err != nil {
		t.Fatalf("backup file not readable: %v", err)
	}
	if string(data) != "SECRET=hello\n" {
		t.Errorf("backup content mismatch: %q", string(data))
	}
	if !strings.Contains(filepath.Base(got), ".bak") {
		t.Errorf("expected .bak suffix in %q", got)
	}
}

func TestBackup_NoSourceFile(t *testing.T) {
	tmpDir := t.TempDir()
	path, err := Backup(filepath.Join(tmpDir, "missing.env"), tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "" {
		t.Errorf("expected empty path for missing source, got %q", path)
	}
}

func TestPruneBackups_RemovesOldFiles(t *testing.T) {
	tmpDir := t.TempDir()

	old := filepath.Join(tmpDir, "old.bak")
	recent := filepath.Join(tmpDir, "recent.bak")

	if err := os.WriteFile(old, []byte("old"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(recent, []byte("new"), 0600); err != nil {
		t.Fatal(err)
	}

	past := time.Now().Add(-48 * time.Hour)
	if err := os.Chtimes(old, past, past); err != nil {
		t.Fatal(err)
	}

	n, err := PruneBackups(tmpDir, 24*time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 removed, got %d", n)
	}
	if _, err := os.Stat(recent); err != nil {
		t.Error("recent file should still exist")
	}
}

func TestPruneBackups_MissingDir(t *testing.T) {
	n, err := PruneBackups("/nonexistent/path", time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0, got %d", n)
	}
}
