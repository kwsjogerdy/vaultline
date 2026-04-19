package envrotate_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/vaultline/vaultline/internal/envrotate"
)

func tmpLedger(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "rotate.json")
}

func TestNew_EmptyWhenMissing(t *testing.T) {
	l, err := envrotate.New(tmpLedger(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := l.Get("MY_KEY")
	if ok {
		t.Error("expected no record for unknown key")
	}
}

func TestRotate_IncrementsVersion(t *testing.T) {
	l, _ := envrotate.New(tmpLedger(t))
	r1 := l.Rotate("DB_PASS")
	r2 := l.Rotate("DB_PASS")
	if r1.Version != 1 {
		t.Errorf("expected version 1, got %d", r1.Version)
	}
	if r2.Version != 2 {
		t.Errorf("expected version 2, got %d", r2.Version)
	}
}

func TestSaveAndReload(t *testing.T) {
	path := tmpLedger(t)
	l, _ := envrotate.New(path)
	l.Rotate("API_KEY")
	if err := l.Save(); err != nil {
		t.Fatalf("save: %v", err)
	}
	l2, err := envrotate.New(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	r, ok := l2.Get("API_KEY")
	if !ok {
		t.Fatal("expected record after reload")
	}
	if r.Version != 1 {
		t.Errorf("expected version 1, got %d", r.Version)
	}
}

func TestDueForRotation(t *testing.T) {
	l, _ := envrotate.New(tmpLedger(t))
	l.Rotate("FRESH_KEY")
	keys := []string{"FRESH_KEY", "OLD_KEY", "NEW_KEY"}
	// OLD_KEY and NEW_KEY have no record — both due
	due := l.DueForRotation(keys, 24*time.Hour)
	if len(due) != 2 {
		t.Errorf("expected 2 due keys, got %d: %v", len(due), due)
	}
}

func TestSave_BadPath(t *testing.T) {
	l, _ := envrotate.New("/nonexistent/dir/rotate.json")
	l.Rotate("X")
	if err := l.Save(); err == nil {
		t.Error("expected error writing to bad path")
	}
	os.Remove("/nonexistent/dir/rotate.json")
}
