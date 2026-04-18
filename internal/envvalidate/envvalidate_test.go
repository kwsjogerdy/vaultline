package envvalidate_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envvalidate"
)

func TestValidate_AllValid(t *testing.T) {
	v := envvalidate.New(map[string]envvalidate.Rule{
		"PORT": {Required: true, Type: "int"},
		"DEBUG": {Required: false, Type: "bool"},
	})
	results := v.Validate(map[string]string{"PORT": "8080", "DEBUG": "true"})
	if len(results) != 0 {
		t.Fatalf("expected no failures, got %v", results)
	}
}

func TestValidate_MissingRequired(t *testing.T) {
	v := envvalidate.New(map[string]envvalidate.Rule{
		"API_KEY": {Required: true},
	})
	results := v.Validate(map[string]string{})
	if len(results) != 1 || results[0].Key != "API_KEY" {
		t.Fatalf("expected missing API_KEY failure, got %v", results)
	}
}

func TestValidate_WrongType_Int(t *testing.T) {
	v := envvalidate.New(map[string]envvalidate.Rule{
		"PORT": {Type: "int"},
	})
	results := v.Validate(map[string]string{"PORT": "not-a-number"})
	if len(results) != 1 {
		t.Fatalf("expected 1 failure, got %v", results)
	}
}

func TestValidate_WrongType_Bool(t *testing.T) {
	v := envvalidate.New(map[string]envvalidate.Rule{
		"ENABLED": {Type: "bool"},
	})
	results := v.Validate(map[string]string{"ENABLED": "yes-please"})
	if len(results) != 1 {
		t.Fatalf("expected 1 failure, got %v", results)
	}
}

func TestValidate_PatternMismatch(t *testing.T) {
	v := envvalidate.New(map[string]envvalidate.Rule{
		"ENV": {Pattern: "^(dev|staging|prod)$"},
	})
	results := v.Validate(map[string]string{"ENV": "local"})
	if len(results) != 1 {
		t.Fatalf("expected pattern failure, got %v", results)
	}
}

func TestValidate_PatternMatch(t *testing.T) {
	v := envvalidate.New(map[string]envvalidate.Rule{
		"ENV": {Pattern: "^(dev|staging|prod)$"},
	})
	results := v.Validate(map[string]string{"ENV": "prod"})
	if len(results) != 0 {
		t.Fatalf("expected no failures, got %v", results)
	}
}

func TestValidate_OptionalMissing_NoError(t *testing.T) {
	v := envvalidate.New(map[string]envvalidate.Rule{
		"OPTIONAL": {Required: false, Type: "int"},
	})
	results := v.Validate(map[string]string{})
	if len(results) != 0 {
		t.Fatalf("expected no failures for optional missing key, got %v", results)
	}
}
