package envdefaults_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/vaultline/vaultline/internal/envdefaults"
)

func TestApply_FillsMissingKeys(t *testing.T) {
	a, _ := envdefaults.New([]envdefaults.Default{
		{Key: "LOG_LEVEL", Value: "info"},
		{Key: "TIMEOUT", Value: "30s"},
	})
	result := a.Apply(map[string]string{"LOG_LEVEL": "debug"})
	if result["LOG_LEVEL"] != "debug" {
		t.Errorf("expected existing key to be preserved")
	}
	if result["TIMEOUT"] != "30s" {
		t.Errorf("expected default TIMEOUT to be applied")
	}
}

func TestApply_EmptySecrets(t *testing.T) {
	a, _ := envdefaults.New([]envdefaults.Default{{Key: "PORT", Value: "8080"}})
	result := a.Apply(map[string]string{})
	if result["PORT"] != "8080" {
		t.Errorf("expected PORT default")
	}
}

func TestMissing_ReturnsAbsentKeys(t *testing.T) {
	a, _ := envdefaults.New([]envdefaults.Default{
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
	})
	missing := a.Missing(map[string]string{"A": "x"})
	if len(missing) != 1 || missing[0] != "B" {
		t.Errorf("expected [B], got %v", missing)
	}
}

func TestNew_EmptyKeyReturnsError(t *testing.T) {
	_, err := envdefaults.New([]envdefaults.Default{{Key: "", Value: "v"}})
	if err == nil {
		t.Error("expected error for empty key")
	}
}

func TestLoadFile_ParsesDefaults(t *testing.T) {
	defaults := []envdefaults.Default{
		{Key: "FOO", Value: "bar", Description: "test"},
	}
	data, _ := json.Marshal(defaults)
	tmp := filepath.Join(t.TempDir(), "defaults.json")
	os.WriteFile(tmp, data, 0644)

	loaded, err := envdefaults.LoadFile(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(loaded) != 1 || loaded[0].Key != "FOO" {
		t.Errorf("unexpected loaded defaults: %v", loaded)
	}
}

func TestLoadFile_MissingFile(t *testing.T) {
	_, err := envdefaults.LoadFile("/nonexistent/defaults.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
