package envtypecheck_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envtypecheck"
)

func TestInfer_Bool(t *testing.T) {
	for _, v := range []string{"true", "false", "yes", "no", "1", "0", "True", "YES"} {
		if got := envtypecheck.Infer(v); got != envtypecheck.TypeBool {
			t.Errorf("Infer(%q) = %s, want bool", v, got)
		}
	}
}

func TestInfer_Int(t *testing.T) {
	for _, v := range []string{"42", "-7", "1000"} {
		if got := envtypecheck.Infer(v); got != envtypecheck.TypeInt {
			t.Errorf("Infer(%q) = %s, want int", v, got)
		}
	}
}

func TestInfer_Float(t *testing.T) {
	for _, v := range []string{"3.14", "-0.5", "1.0e10"} {
		if got := envtypecheck.Infer(v); got != envtypecheck.TypeFloat {
			t.Errorf("Infer(%q) = %s, want float", v, got)
		}
	}
}

func TestInfer_URL(t *testing.T) {
	for _, v := range []string{"https://example.com", "http://localhost:8200"} {
		if got := envtypecheck.Infer(v); got != envtypecheck.TypeURL {
			t.Errorf("Infer(%q) = %s, want url", v, got)
		}
	}
}

func TestInfer_Email(t *testing.T) {
	if got := envtypecheck.Infer("user@example.com"); got != envtypecheck.TypeEmail {
		t.Errorf("got %s, want email", got)
	}
}

func TestInfer_String(t *testing.T) {
	for _, v := range []string{"hello", "some-secret-value", ""} {
		if got := envtypecheck.Infer(v); got != envtypecheck.TypeString {
			t.Errorf("Infer(%q) = %s, want string", v, got)
		}
	}
}

func TestCheck_NoEnforce_NoErrors(t *testing.T) {
	c := envtypecheck.New(nil)
	secrets := map[string]string{
		"PORT":     "8080",
		"DEBUG":    "true",
		"APP_NAME": "vaultline",
	}
	results := c.Check(secrets)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if envtypecheck.HasErrors(results) {
		t.Error("expected no errors without enforcement")
	}
}

func TestCheck_EnforceMismatch_ReturnsError(t *testing.T) {
	c := envtypecheck.New(map[string]envtypecheck.Type{
		"PORT": envtypecheck.TypeString,
	})
	results := c.Check(map[string]string{"PORT": "8080"})
	if !envtypecheck.HasErrors(results) {
		t.Error("expected type mismatch error")
	}
	if results[0].Error == nil {
		t.Error("expected error on PORT result")
	}
}

func TestCheck_EnforceMatch_NoError(t *testing.T) {
	c := envtypecheck.New(map[string]envtypecheck.Type{
		"PORT": envtypecheck.TypeInt,
	})
	results := c.Check(map[string]string{"PORT": "9000"})
	if envtypecheck.HasErrors(results) {
		t.Errorf("unexpected error: %v", results[0].Error)
	}
}

func TestHasErrors_False(t *testing.T) {
	results := []envtypecheck.Result{{Key: "K", Type: envtypecheck.TypeString}}
	if envtypecheck.HasErrors(results) {
		t.Error("expected no errors")
	}
}
