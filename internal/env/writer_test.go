package env

import (
	"os"
	"strings"
	"testing"
)

func TestWrite_CreatesFile(t *testing.T) {
	tmp, err := os.CreateTemp(t.TempDir(), ".env")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	w := NewWriter(tmp.Name())
	secrets := map[string]string{
		"db_password": "s3cr3t",
		"api_key":     "abc123",
	}

	if err := w.Write(secrets); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	content, err := os.ReadFile(tmp.Name())
	if err != nil {
		t.Fatal(err)
	}

	for _, key := range []string{"DB_PASSWORD", "API_KEY"} {
		if !strings.Contains(string(content), key) {
			t.Errorf("expected key %q in output, got:\n%s", key, content)
		}
	}
}

func TestWrite_QuotesValuesWithSpaces(t *testing.T) {
	tmp, err := os.CreateTemp(t.TempDir(), ".env")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	w := NewWriter(tmp.Name())
	secrets := map[string]string{
		"greeting": "hello world",
	}

	if err := w.Write(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(tmp.Name())
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(content), `GREETING="hello world"`) {
		t.Errorf("expected quoted value, got:\n%s", content)
	}
}

func TestWrite_InvalidPath(t *testing.T) {
	w := NewWriter("/nonexistent/path/.env")
	err := w.Write(map[string]string{"key": "value"})
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}
