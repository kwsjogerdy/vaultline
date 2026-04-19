package envclone_test

import (
	"errors"
	"testing"

	"github.com/vaultline/vaultline/internal/envclone"
)

type mockFetcher struct {
	data map[string]string
	err  error
}

func (m *mockFetcher) GetSecrets(_ string) (map[string]string, error) {
	return m.data, m.err
}

type mockWriter struct {
	written map[string]string
	err     error
}

func (m *mockWriter) PutSecrets(_ string, secrets map[string]string) error {
	m.written = secrets
	return m.err
}

func TestClone_CopiesAllSecrets(t *testing.T) {
	f := &mockFetcher{data: map[string]string{"A": "1", "B": "2"}}
	w := &mockWriter{}
	c := envclone.New(f, w)

	n, err := c.Clone("secret/src", "secret/dst", envclone.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 keys cloned, got %d", n)
	}
	if w.written["A"] != "1" || w.written["B"] != "2" {
		t.Errorf("unexpected written secrets: %v", w.written)
	}
}

func TestClone_OnlyFilter(t *testing.T) {
	f := &mockFetcher{data: map[string]string{"A": "1", "B": "2", "C": "3"}}
	w := &mockWriter{}
	c := envclone.New(f, w)

	n, err := c.Clone("secret/src", "secret/dst", envclone.Options{Only: []string{"A", "C"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2, got %d", n)
	}
	if _, ok := w.written["B"]; ok {
		t.Error("B should have been excluded")
	}
}

func TestClone_RemapsKeys(t *testing.T) {
	f := &mockFetcher{data: map[string]string{"OLD_KEY": "val"}}
	w := &mockWriter{}
	c := envclone.New(f, w)

	_, err := c.Clone("src", "dst", envclone.Options{Remap: map[string]string{"OLD_KEY": "NEW_KEY"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.written["NEW_KEY"] != "val" {
		t.Errorf("expected NEW_KEY=val, got %v", w.written)
	}
}

func TestClone_FetchError(t *testing.T) {
	f := &mockFetcher{err: errors.New("vault down")}
	w := &mockWriter{}
	c := envclone.New(f, w)

	_, err := c.Clone("src", "dst", envclone.Options{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestClone_EmptyPath(t *testing.T) {
	c := envclone.New(&mockFetcher{}, &mockWriter{})
	_, err := c.Clone("", "dst", envclone.Options{})
	if err == nil {
		t.Fatal("expected error for empty src")
	}
}
