package envduration_test

import (
	"testing"
	"time"

	"github.com/vaultline/vaultline/internal/envduration"
)

func TestParse_AllKeys_SkipsUnparseable(t *testing.T) {
	secrets := map[string]string{
		"TIMEOUT":  "30s",
		"INTERVAL": "5m",
		"NAME":     "alice", // not a duration – should be skipped
	}
	p := envduration.New()
	results, err := p.Parse(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestParse_TargetedKeys_ReturnsError(t *testing.T) {
	secrets := map[string]string{
		"TIMEOUT": "not-a-duration",
	}
	p := envduration.New("TIMEOUT")
	_, err := p.Parse(secrets)
	if err == nil {
		t.Fatal("expected error for unparseable targeted key")
	}
}

func TestParse_TargetedKeys_Success(t *testing.T) {
	secrets := map[string]string{
		"TIMEOUT":  "2h",
		"UNRELATED": "hello",
	}
	p := envduration.New("TIMEOUT")
	results, err := p.Parse(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Duration != 2*time.Hour {
		t.Errorf("expected 2h, got %v", results[0].Duration)
	}
}

func TestParse_EmptySecrets(t *testing.T) {
	p := envduration.New()
	results, err := p.Parse(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected empty results")
	}
}

func TestFormat_ReplacesKnownKeys(t *testing.T) {
	secrets := map[string]string{
		"TIMEOUT": "30s",
		"NAME":    "bob",
	}
	out := envduration.Format(secrets, "TIMEOUT")
	if out["TIMEOUT"] != "30s" {
		t.Errorf("expected '30s', got %q", out["TIMEOUT"])
	}
	if out["NAME"] != "bob" {
		t.Errorf("expected 'bob', got %q", out["NAME"])
	}
}

func TestFormat_SkipsUnparseable(t *testing.T) {
	secrets := map[string]string{
		"TIMEOUT": "not-a-duration",
	}
	out := envduration.Format(secrets, "TIMEOUT")
	// value unchanged when it cannot be re-formatted
	if out["TIMEOUT"] != "not-a-duration" {
		t.Errorf("expected original value preserved, got %q", out["TIMEOUT"])
	}
}
