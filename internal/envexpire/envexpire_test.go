package envexpire_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/vaultline/vaultline/internal/envexpire"
)

func tmpFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "expire.json")
}

func TestSet_IsExpired_NotYet(t *testing.T) {
	s, _ := envexpire.New(tmpFile(t))
	s.Set("MY_KEY", 10*time.Minute)
	if s.IsExpired("MY_KEY") {
		t.Fatal("expected key not to be expired")
	}
}

func TestIsExpired_AfterTTL(t *testing.T) {
	s, _ := envexpire.New(tmpFile(t))
	s.Set("OLD_KEY", -1*time.Second)
	if !s.IsExpired("OLD_KEY") {
		t.Fatal("expected key to be expired")
	}
}

func TestIsExpired_UnknownKey(t *testing.T) {
	s, _ := envexpire.New(tmpFile(t))
	if s.IsExpired("GHOST") {
		t.Fatal("unknown key should not be expired")
	}
}

func TestExpired_ReturnsExpiredKeys(t *testing.T) {
	s, _ := envexpire.New(tmpFile(t))
	s.Set("FRESH", 10*time.Minute)
	s.Set("STALE", -1*time.Second)
	exp := s.Expired()
	if len(exp) != 1 || exp[0] != "STALE" {
		t.Fatalf("expected [STALE], got %v", exp)
	}
}

func TestSaveAndReload(t *testing.T) {
	path := tmpFile(t)
	s, _ := envexpire.New(path)
	s.Set("TOKEN", 5*time.Minute)
	if err := s.Save(); err != nil {
		t.Fatalf("save: %v", err)
	}
	s2, err := envexpire.New(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if s2.IsExpired("TOKEN") {
		t.Fatal("reloaded key should not be expired")
	}
}

func TestNew_MissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")
	s, err := envexpire.New(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Expired()) != 0 {
		t.Fatal("expected empty store")
	}
}

func TestNew_CorruptFile(t *testing.T) {
	path := tmpFile(t)
	os.WriteFile(path, []byte("not json"), 0600)
	_, err := envexpire.New(path)
	if err == nil {
		t.Fatal("expected parse error")
	}
}
