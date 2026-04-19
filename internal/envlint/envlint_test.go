package envlint_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envlint"
)

func TestLint_ValidSecrets(t *testing.T) {
	l := envlint.New(false)
	secrets := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"API_KEY":      "abc123",
	}
	findings := l.Lint(secrets)
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d: %v", len(findings), findings)
	}
}

func TestLint_InvalidKeyFormat(t *testing.T) {
	l := envlint.New(false)
	secrets := map[string]string{
		"bad-key": "value",
	}
	findings := l.Lint(secrets)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Rule != "key-format" {
		t.Errorf("expected key-format rule, got %s", findings[0].Rule)
	}
	if findings[0].Severity != envlint.Error {
		t.Errorf("expected ERROR severity")
	}
}

func TestLint_EmptyValue(t *testing.T) {
	l := envlint.New(false)
	secrets := map[string]string{
		"MY_KEY": "",
	}
	findings := l.Lint(secrets)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Rule != "empty-value" {
		t.Errorf("expected empty-value rule, got %s", findings[0].Rule)
	}
	if findings[0].Severity != envlint.Warn {
		t.Errorf("expected WARN severity")
	}
}

func TestLint_AllowEmpty(t *testing.T) {
	l := envlint.New(true)
	secrets := map[string]string{"MY_KEY": ""}
	findings := l.Lint(secrets)
	if len(findings) != 0 {
		t.Fatalf("expected no findings when allowEmpty=true, got %d", len(findings))
	}
}

func TestLint_MultilineValue(t *testing.T) {
	l := envlint.New(false)
	secrets := map[string]string{"CERT": "line1\nline2"}
	findings := l.Lint(secrets)
	if len(findings) != 1 || findings[0].Rule != "multiline-value" {
		t.Errorf("expected multiline-value finding")
	}
}

func TestHasErrors_True(t *testing.T) {
	findings := []envlint.Finding{
		{Rule: "key-format", Severity: envlint.Error},
	}
	if !envlint.HasErrors(findings) {
		t.Error("expected HasErrors to return true")
	}
}

func TestHasErrors_False(t *testing.T) {
	findings := []envlint.Finding{
		{Rule: "empty-value", Severity: envlint.Warn},
	}
	if envlint.HasErrors(findings) {
		t.Error("expected HasErrors to return false")
	}
}
