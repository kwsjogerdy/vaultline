package envredact_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envredact"
)

func TestApply_RedactsSensitiveKeys(t *testing.T) {
	r := envredact.New()
	secrets := map[string]string{
		"DB_PASSWORD": "supersecret",
		"APP_NAME":    "myapp",
		"API_KEY":     "abc123",
		"HOST":        "localhost",
	}
	res := r.Apply(secrets)
	if res.Secrets["DB_PASSWORD"] != "[REDACTED]" {
		t.Errorf("expected DB_PASSWORD to be redacted, got %q", res.Secrets["DB_PASSWORD"])
	}
	if res.Secrets["API_KEY"] != "[REDACTED]" {
		t.Errorf("expected API_KEY to be redacted, got %q", res.Secrets["API_KEY"])
	}
	if res.Secrets["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME to be unchanged, got %q", res.Secrets["APP_NAME"])
	}
	if res.Secrets["HOST"] != "localhost" {
		t.Errorf("expected HOST to be unchanged")
	}
}

func TestApply_ReportsRedactedKeys(t *testing.T) {
	r := envredact.New()
	secrets := map[string]string{
		"SECRET_KEY": "val1",
		"SAFE_KEY":   "val2",
	}
	res := r.Apply(secrets)
	if len(res.Redacted) != 1 || res.Redacted[0] != "SECRET_KEY" {
		t.Errorf("expected [SECRET_KEY] redacted, got %v", res.Redacted)
	}
}

func TestApply_EmptySecrets(t *testing.T) {
	r := envredact.New()
	res := r.Apply(map[string]string{})
	if len(res.Secrets) != 0 {
		t.Errorf("expected empty result")
	}
	if len(res.Redacted) != 0 {
		t.Errorf("expected no redacted keys")
	}
}

func TestApply_CustomPlaceholder(t *testing.T) {
	r := envredact.NewWithPatterns([]string{`(?i)token`}, "***")
	secrets := map[string]string{"AUTH_TOKEN": "xyz", "NAME": "foo"}
	res := r.Apply(secrets)
	if res.Secrets["AUTH_TOKEN"] != "***" {
		t.Errorf("expected *** placeholder, got %q", res.Secrets["AUTH_TOKEN"])
	}
	if res.Secrets["NAME"] != "foo" {
		t.Errorf("expected NAME unchanged")
	}
}

func TestIsSensitive_MatchesPatterns(t *testing.T) {
	r := envredact.New()
	cases := []struct {
		key       string
		expected  bool
	}{
		{"DB_PASSWORD", true},
		{"api_key", true},
		{"PRIVATE_KEY", true},
		{"APP_ENV", false},
		{"PORT", false},
		{"AUTH_HEADER", true},
	}
	for _, tc := range cases {
		got := r.IsSensitive(tc.key)
		if got != tc.expected {
			t.Errorf("IsSensitive(%q) = %v, want %v", tc.key, got, tc.expected)
		}
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	r := envredact.New()
	orig := map[string]string{"PASSWORD": "secret123", "HOST": "db"}
	r.Apply(orig)
	if orig["PASSWORD"] != "secret123" {
		t.Error("Apply must not mutate the input map")
	}
}
