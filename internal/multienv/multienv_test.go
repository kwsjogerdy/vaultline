package multienv_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vaultline/vaultline/internal/multienv"
)

func TestWriteAll_AllTargetsReceiveSecrets(t *testing.T) {
	dir := t.TempDir()
	p1 := filepath.Join(dir, "a.env")
	p2 := filepath.Join(dir, "b.env")

	w := multienv.New([]multienv.Target{
		{Path: p1},
		{Path: p2},
	})

	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if err := w.WriteAll(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, p := range []string{p1, p2} {
		b, err := os.ReadFile(p)
		if err != nil {
			t.Fatalf("read %s: %v", p, err)
		}
		if !strings.Contains(string(b), "FOO") {
			t.Errorf("%s missing FOO", p)
		}
	}
}

func TestWriteAll_FiltersByKeys(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "filtered.env")

	w := multienv.New([]multienv.Target{
		{Path: p, Keys: []string{"ONLY"}},
	})

	secrets := map[string]string{"ONLY": "yes", "SKIP": "no"}
	if err := w.WriteAll(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	b, _ := os.ReadFile(p)
	content := string(b)
	if !strings.Contains(content, "ONLY") {
		t.Error("expected ONLY in output")
	}
	if strings.Contains(content, "SKIP") {
		t.Error("SKIP should not appear in filtered output")
	}
}

func TestWriteAll_InvalidPath(t *testing.T) {
	w := multienv.New([]multienv.Target{
		{Path: "/no/such/dir/out.env"},
	})
	if err := w.WriteAll(map[string]string{"K": "v"}); err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestWriteAll_EmptyTargets(t *testing.T) {
	w := multienv.New(nil)
	if err := w.WriteAll(map[string]string{"K": "v"}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
