package envsnapshot_test

import (
	"os"
	"testing"
	"time"

	"github.com/vaultline/vaultline/internal/envsnapshot"
)

func tmpDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "snap-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func TestSave_AndList(t *testing.T) {
	m := envsnapshot.New(tmpDir(t))
	secrets := map[string]string{"DB_URL": "postgres://localhost", "API_KEY": "abc"}
	if err := m.Save("dev", secrets); err != nil {
		t.Fatalf("Save: %v", err)
	}
	list, err := m.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 snapshot, got %d", len(list))
	}
	if list[0].Label != "dev" {
		t.Errorf("label = %q, want dev", list[0].Label)
	}
	if list[0].Secrets["DB_URL"] != "postgres://localhost" {
		t.Errorf("DB_URL mismatch")
	}
}

func TestList_EmptyDir(t *testing.T) {
	m := envsnapshot.New(tmpDir(t))
	list, err := m.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 0 {
		t.Errorf("expected empty list")
	}
}

func TestList_NewestFirst(t *testing.T) {
	m := envsnapshot.New(tmpDir(t))
	m.Save("first", map[string]string{"A": "1"})
	time.Sleep(2 * time.Millisecond)
	m.Save("second", map[string]string{"B": "2"})
	list, err := m.List()
	if err != nil {
		t.Fatal(err)
	}
	if list[0].Label != "second" {
		t.Errorf("expected newest first, got %q", list[0].Label)
	}
}

func TestLatest_ReturnsNewest(t *testing.T) {
	m := envsnapshot.New(tmpDir(t))
	m.Save("old", map[string]string{})
	time.Sleep(2 * time.Millisecond)
	m.Save("new", map[string]string{"X": "y"})
	snap, err := m.Latest()
	if err != nil {
		t.Fatal(err)
	}
	if snap.Label != "new" {
		t.Errorf("latest label = %q, want new", snap.Label)
	}
}

func TestLatest_MissingDir(t *testing.T) {
	m := envsnapshot.New("/tmp/no-such-snap-dir-xyz")
	snap, err := m.Latest()
	if err != nil {
		t.Fatal(err)
	}
	if snap != nil {
		t.Errorf("expected nil snapshot for missing dir")
	}
}
