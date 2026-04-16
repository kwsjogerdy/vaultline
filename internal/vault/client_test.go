package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultline/internal/vault"
)

func makeVaultResponse(data map[string]interface{}) []byte {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"data": data,
		},
	}
	b, _ := json.Marshal(payload)
	return b
}

func TestGetSecrets_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") != "test-token" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(makeVaultResponse(map[string]interface{}{"API_KEY": "abc123", "DB_PASS": "secret"}))
	}))
	defer ts.Close()

	client := vault.NewClient(ts.URL, "test-token")
	secrets, err := client.GetSecrets("secret", "myapp")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if secrets["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %s", secrets["API_KEY"])
	}
	if secrets["DB_PASS"] != "secret" {
		t.Errorf("expected DB_PASS=secret, got %s", secrets["DB_PASS"])
	}
}

func TestGetSecrets_Unauthorized(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	client := vault.NewClient(ts.URL, "bad-token")
	_, err := client.GetSecrets("secret", "myapp")
	if err == nil {
		t.Fatal("expected error for unauthorized request")
	}
}

func TestGetSecrets_NotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	client := vault.NewClient(ts.URL, "test-token")
	_, err := client.GetSecrets("secret", "nonexistent")
	if err == nil {
		t.Fatal("expected error for not found secret")
	}
}
