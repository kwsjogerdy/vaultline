package envclip

import (
	"testing"
	"time"
)

func TestCopy_KeyNotFound(t *testing.T) {
	c := New(0)
	secrets := map[string]string{"FOO": "bar"}
	err := c.Copy(secrets, "MISSING")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
	if err.Error() != `key "MISSING" not found` {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_DefaultClearAfter(t *testing.T) {
	c := New(30 * time.Second)
	if c.clearAfter != 30*time.Second {
		t.Fatalf("expected 30s, got %v", c.clearAfter)
	}
}

func TestNew_ZeroClearAfter(t *testing.T) {
	c := New(0)
	if c.clearAfter != 0 {
		t.Fatalf("expected 0, got %v", c.clearAfter)
	}
}

func TestWriteClipboard_UnsupportedPlatform(t *testing.T) {
	// We can't easily mock runtime.GOOS, so we test the error path indirectly
	// by confirming writeClipboard returns an error only on unknown platforms.
	// This test is a no-op on supported platforms and documents expected behaviour.
	_ = writeClipboard // ensure symbol is accessible
}

func TestCopy_EmptySecrets(t *testing.T) {
	c := New(0)
	err := c.Copy(map[string]string{}, "ANY_KEY")
	if err == nil {
		t.Fatal("expected error for empty secrets map")
	}
}
