package schema_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/schema"
)

func TestValidate_AllPresent(t *testing.T) {
	s := schema.New([]schema.Rule{
		{Key: "DB_URL", Required: true},
		{Key: "API_KEY", Required: true},
	})
	err := s.Validate(map[string]string{"DB_URL": "postgres://", "API_KEY": "abc"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidate_MissingRequired(t *testing.T) {
	s := schema.New([]schema.Rule{
		{Key: "DB_URL", Required: true},
	})
	err := s.Validate(map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing required key")
	}
	ve, ok := err.(*schema.ValidationError)
	if !ok || len(ve.Violations) != 1 {
		t.Fatalf("unexpected error type or count: %v", err)
	}
}

func TestValidate_PatternMismatch(t *testing.T) {
	s := schema.New([]schema.Rule{
		{Key: "DB_URL", Required: true, Pattern: "postgres"},
	})
	err := s.Validate(map[string]string{"DB_URL": "mysql://localhost"})
	if err == nil {
		t.Fatal("expected pattern mismatch error")
	}
}

func TestValidate_PatternMatch(t *testing.T) {
	s := schema.New([]schema.Rule{
		{Key: "DB_URL", Required: true, Pattern: "postgres"},
	})
	err := s.Validate(map[string]string{"DB_URL": "postgres://localhost"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidate_OptionalMissing(t *testing.T) {
	s := schema.New([]schema.Rule{
		{Key: "OPTIONAL_KEY", Required: false},
	})
	err := s.Validate(map[string]string{})
	if err != nil {
		t.Fatalf("optional missing key should not error, got %v", err)
	}
}

func TestValidate_MultipleViolations(t *testing.T) {
	s := schema.New([]schema.Rule{
		{Key: "A", Required: true},
		{Key: "B", Required: true},
	})
	err := s.Validate(map[string]string{})
	ve, ok := err.(*schema.ValidationError)
	if !ok || len(ve.Violations) != 2 {
		t.Fatalf("expected 2 violations, got %v", err)
	}
}
