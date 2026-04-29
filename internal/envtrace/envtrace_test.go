package envtrace_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/vaultline/vaultline/internal/envtrace"
)

func tmpFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "trace.json")
}

func TestRecord_AddsEntry(t *testing.T) {
	tr, err := envtrace.New(tmpFile(t))
	if err != nil {
		t.Fatal(err)
	}
	tr.Record("DB_PASSWORD", "sync")
	entries := tr.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Key != "DB_PASSWORD" || entries[0].Operation != "sync" {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
}

func TestForKey_FiltersCorrectly(t *testing.T) {
	tr, _ := envtrace.New(tmpFile(t))
	tr.Record("KEY_A", "read")
	tr.Record("KEY_B", "sync")
	tr.Record("KEY_A", "export")

	results := tr.ForKey("KEY_A")
	if len(results) != 2 {
		t.Fatalf("expected 2, got %d", len(results))
	}
	for _, e := range results {
		if e.Key != "KEY_A" {
			t.Errorf("unexpected key: %s", e.Key)
		}
	}
}

func TestSave_AndReload(t *testing.T) {
	path := tmpFile(t)
	tr, _ := envtrace.New(path)
	tr.Record("SECRET", "import")
	if err := tr.Save(); err != nil {
		t.Fatal(err)
	}

	tr2, err := envtrace.New(path)
	if err != nil {
		t.Fatal(err)
	}
	entries := tr2.Entries()
	if len(entries) != 1 || entries[0].Key != "SECRET" {
		t.Errorf("reload failed: %+v", entries)
	}
}

func TestEntries_SortedByTime(t *testing.T) {
	tr, _ := envtrace.New(tmpFile(t))
	now := time.Now()
	// Inject out-of-order entries via Save/Reload trick using raw JSON.
	raw := []envtrace.Entry{
		{Key: "B", Operation: "op", At: now.Add(2 * time.Second)},
		{Key: "A", Operation: "op", At: now.Add(1 * time.Second)},
	}
	path := tmpFile(t)
	data, _ := json.Marshal(raw)
	os.WriteFile(path, data, 0600)

	tr2, _ := envtrace.New(path)
	entries := tr2.Entries()
	if entries[0].Key != "A" || entries[1].Key != "B" {
		t.Errorf("expected sorted order A,B got %s,%s", entries[0].Key, entries[1].Key)
	}
}

func TestNew_CorruptFile_ReturnsError(t *testing.T) {
	path := tmpFile(t)
	os.WriteFile(path, []byte("not json"), 0600)
	_, err := envtrace.New(path)
	if err == nil {
		t.Fatal("expected error for corrupt file")
	}
}
