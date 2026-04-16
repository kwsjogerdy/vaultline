package sync_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/vaultline/internal/env"
	"github.com/vaultline/internal/sync"
	"github.com/vaultline/internal/vault"
)

func makeServer(t *testing.T, data map[string]string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"data": data},
		})
	}))
}

func TestSync_WritesSecrets(t *testing.T) {
	server := makeServer(t, map[string]string{"DB_PASS": "secret", "API_KEY": "abc123"})
	defer server.Close()

	client, err := vault.NewClient(server.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	dir := t.TempDir()
	outPath := filepath.Join(dir, ".env")
	writer := env.NewWriter()
	syncer := sync.New(client, writer)

	result, err := syncer.Sync("secret/data/app", outPath, nil)
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}

	if result.Count != 2 {
		t.Errorf("expected 2 secrets written, got %d", result.Count)
	}
	if result.Skipped != 0 {
		t.Errorf("expected 0 skipped, got %d", result.Skipped)
	}
	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		t.Error("expected output file to exist")
	}
}

func TestSync_SkipsKeys(t *testing.T) {
	server := makeServer(t, map[string]string{"DB_PASS": "secret", "API_KEY": "abc123", "DEBUG": "true"})
	defer server.Close()

	client, err := vault.NewClient(server.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	dir := t.TempDir()
	writer := env.NewWriter()
	syncer := sync.New(client, writer)

	result, err := syncer.Sync("secret/data/app", filepath.Join(dir, ".env"), []string{"DEBUG"})
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}

	if result.Count != 2 {
		t.Errorf("expected 2 secrets written, got %d", result.Count)
	}
	if result.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", result.Skipped)
	}
}
