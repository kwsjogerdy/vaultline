package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_MissingToken(t *testing.T) {
	os.Unsetenv("VAULTLINE_VAULT_TOKEN")
	_, err := Load("")
	if err == nil {
		t.Fatal("expected error when vault_token is missing, got nil")
	}
}

func TestLoad_FromEnvVars(t *testing.T) {
	t.Setenv("VAULTLINE_VAULT_ADDR", "http://localhost:8200")
	t.Setenv("VAULTLINE_VAULT_TOKEN", "test-token")
	t.Setenv("VAULTLINE_SECRET_PATH", "secret/data/myapp")
	t.Setenv("VAULTLINE_ENV_FILE", ".env.local")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.VaultToken != "test-token" {
		t.Errorf("expected token 'test-token', got '%s'", cfg.VaultToken)
	}
	if cfg.SecretPath != "secret/data/myapp" {
		t.Errorf("expected secret path 'secret/data/myapp', got '%s'", cfg.SecretPath)
	}
	if cfg.EnvFile != ".env.local" {
		t.Errorf("expected env file '.env.local', got '%s'", cfg.EnvFile)
	}
}

func TestLoad_FromFile(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")
	content := []byte("vault_addr: http://vault:8200\nvault_token: file-token\nsecret_path: secret/data/app\nenv_file: .env\n")
	if err := os.WriteFile(cfgPath, content, 0600); err != nil {
		t.Fatal(err)
	}

	os.Unsetenv("VAULTLINE_VAULT_TOKEN")

	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.VaultToken != "file-token" {
		t.Errorf("expected 'file-token', got '%s'", cfg.VaultToken)
	}
	if cfg.VaultAddr != "http://vault:8200" {
		t.Errorf("expected 'http://vault:8200', got '%s'", cfg.VaultAddr)
	}
}
