package envrenum_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envrenum"
)

func TestNew_EmptyKey_ReturnsError(t *testing.T) {
	_, err := envrenum.New([]envrenum.Rule{{Key: "", Allowed: []string{"a"}}})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestNew_NoAllowedValues_ReturnsError(t *testing.T) {
	_, err := envrenum.New([]envrenum.Rule{{Key: "ENV", Allowed: nil}})
	if err == nil {
		t.Fatal("expected error for empty allowed list")
	}
}

func TestValidate_AllValid(t *testing.T) {
	v, err := envrenum.New([]envrenum.Rule{
		{Key: "ENV", Allowed: []string{"dev", "staging", "prod"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	secrets := map[string]string{"ENV": "dev"}
	violations := v.Validate(secrets)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestValidate_InvalidValue(t *testing.T) {
	v, err := envrenum.New([]envrenum.Rule{
		{Key: "LOG_LEVEL", Allowed: []string{"debug", "info", "warn", "error"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	secrets := map[string]string{"LOG_LEVEL": "verbose"}
	violations := v.Validate(secrets)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "LOG_LEVEL" {
		t.Errorf("unexpected key %q", violations[0].Key)
	}
	if violations[0].Value != "verbose" {
		t.Errorf("unexpected value %q", violations[0].Value)
	}
}

func TestValidate_MissingKeySkipped(t *testing.T) {
	v, err := envrenum.New([]envrenum.Rule{
		{Key: "ENV", Allowed: []string{"dev", "prod"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	// ENV not present in secrets — should not produce a violation
	violations := v.Validate(map[string]string{"OTHER": "x"})
	if len(violations) != 0 {
		t.Fatalf("expected no violations for missing key, got %v", violations)
	}
}

func TestValidate_MultipleViolations(t *testing.T) {
	v, err := envrenum.New([]envrenum.Rule{
		{Key: "ENV", Allowed: []string{"dev", "prod"}},
		{Key: "REGION", Allowed: []string{"us-east-1", "eu-west-1"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	secrets := map[string]string{"ENV": "test", "REGION": "ap-south-1"}
	violations := v.Validate(secrets)
	if len(violations) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(violations))
	}
}

func TestHasViolations_False(t *testing.T) {
	v, _ := envrenum.New([]envrenum.Rule{
		{Key: "MODE", Allowed: []string{"read", "write"}},
	})
	if v.HasViolations(map[string]string{"MODE": "read"}) {
		t.Error("expected no violations")
	}
}

func TestViolation_Error_Format(t *testing.T) {
	violation := envrenum.Violation{
		Key:     "ENV",
		Value:   "bad",
		Allowed: []string{"dev", "prod"},
	}
	msg := violation.Error()
	if msg == "" {
		t.Error("expected non-empty error message")
	}
	for _, want := range []string{"ENV", "bad", "dev", "prod"} {
		if !containsStr(msg, want) {
			t.Errorf("expected %q in error message %q", want, msg)
		}
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsSubstr(s, sub))
}

func containsSubstr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
