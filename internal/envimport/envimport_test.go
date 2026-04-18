package envimport_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/vaultline/vaultline/internal/envimport"
)

func writeTmpEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestFromFile_ImportsSecrets(t *testing.T) {
	received := map[string]string{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		for k, v := range body["data"] {
			received[k] = v
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	envFile := writeTmpEnv(t, "DB_HOST=localhost\nDB_PASS=secret\n")
	im := envimport.New(ts.URL, "test-token")
	res, err := im.FromFile(envFile, "secret/data/myapp", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Imported) != 2 {
		t.Errorf("expected 2 imported, got %d", len(res.Imported))
	}
	if len(res.Errors) != 0 {
		t.Errorf("unexpected errors: %v", res.Errors)
	}
}

func TestFromFile_SkipsKeys(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	envFile := writeTmpEnv(t, "API_KEY=abc\nDEBUG=true\n")
	im := envimport.New(ts.URL, "tok")
	res, err := im.FromFile(envFile, "secret/data/app", []string{"DEBUG"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "DEBUG" {
		t.Errorf("expected DEBUG skipped, got %v", res.Skipped)
	}
	if len(res.Imported) != 1 {
		t.Errorf("expected 1 imported, got %d", len(res.Imported))
	}
}

func TestFromFile_MissingFile(t *testing.T) {
	im := envimport.New("http://localhost", "tok")
	_, err := im.FromFile(filepath.Join(t.TempDir(), "missing.env"), "secret/data/app", nil)
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestFromFile_VaultError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	envFile := writeTmpEnv(t, "KEY=val\n")
	im := envimport.New(ts.URL, "bad-token")
	res, err := im.FromFile(envFile, "secret/data/app", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Errors) == 0 {
		t.Error("expected errors for vault 403")
	}
}
