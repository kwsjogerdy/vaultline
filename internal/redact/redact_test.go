package redact_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/redact"
)

func TestIsSensitive_MatchesDefaultPatterns(t *testing.T) {
	r := redact.New()
	sensitive := []string{"DB_PASSWORD", "API_TOKEN", "AWS_SECRET", "PRIVATE_KEY", "AUTH_HEADER"}
	for _, k := range sensitive {
		if !r.IsSensitive(k) {
			t.Errorf("expected %q to be sensitive", k)
		}
	}
}

func TestIsSensitive_SafeKeys(t *testing.T) {
	r := redact.New()
	safe := []string{"APP_ENV", "PORT", "LOG_LEVEL", "REGION"}
	for _, k := range safe {
		if r.IsSensitive(k) {
			t.Errorf("expected %q to NOT be sensitive", k)
		}
	}
}

func TestRedactMap_ReplacesOnlySensitive(t *testing.T) {
	r := redact.New()
	input := map[string]string{
		"DB_PASSWORD": "s3cr3t",
		"APP_ENV":     "production",
		"API_TOKEN":   "tok_abc123",
		"PORT":        "8080",
	}
	out := r.RedactMap(input)

	if out["DB_PASSWORD"] != "[REDACTED]" {
		t.Errorf("expected DB_PASSWORD to be redacted, got %q", out["DB_PASSWORD"])
	}
	if out["API_TOKEN"] != "[REDACTED]" {
		t.Errorf("expected API_TOKEN to be redacted, got %q", out["API_TOKEN"])
	}
	if out["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV to be unchanged, got %q", out["APP_ENV"])
	}
	if out["PORT"] != "8080" {
		t.Errorf("expected PORT to be unchanged, got %q", out["PORT"])
	}
}

func TestRedactValue_SensitiveKey(t *testing.T) {
	r := redact.New()
	got := r.RedactValue("DB_SECRET", "mysecret")
	if got != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", got)
	}
}

func TestRedactValue_SafeKey(t *testing.T) {
	r := redact.New()
	got := r.RedactValue("LOG_LEVEL", "debug")
	if got != "debug" {
		t.Errorf("expected 'debug', got %q", got)
	}
}

func TestNewWithPatterns_CustomPlaceholder(t *testing.T) {
	r := redact.NewWithPatterns([]string{"internal"}, "***")
	if !r.IsSensitive("INTERNAL_KEY") {
		t.Error("expected INTERNAL_KEY to match custom pattern")
	}
	got := r.RedactValue("INTERNAL_KEY", "value")
	if got != "***" {
		t.Errorf("expected ***, got %q", got)
	}
}
