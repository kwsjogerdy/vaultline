package envsanitize_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envsanitize"
)

func TestApply_TrimSpace(t *testing.T) {
	s, err := envsanitize.New(envsanitize.Options{TrimSpace: true})
	if err != nil {
		t.Fatal(err)
	}
	out := s.Apply(map[string]string{"KEY": "  hello  "})
	if out["KEY"] != "hello" {
		t.Errorf("expected 'hello', got %q", out["KEY"])
	}
}

func TestApply_StripNewlines(t *testing.T) {
	s, err := envsanitize.New(envsanitize.Options{StripNewlines: true})
	if err != nil {
		t.Fatal(err)
	}
	out := s.Apply(map[string]string{"KEY": "line1\nline2"})
	if out["KEY"] != "line1 line2" {
		t.Errorf("unexpected value: %q", out["KEY"])
	}
}

func TestApply_StripNulls(t *testing.T) {
	s, err := envsanitize.New(envsanitize.Options{StripNulls: true})
	if err != nil {
		t.Fatal(err)
	}
	out := s.Apply(map[string]string{"KEY": "val\x00ue"})
	if out["KEY"] != "value" {
		t.Errorf("expected 'value', got %q", out["KEY"])
	}
}

func TestApply_CustomRule_Replacement(t *testing.T) {
	s, err := envsanitize.New(envsanitize.Options{
		Rules: []envsanitize.Rule{
			{Pattern: `[^a-zA-Z0-9]`, Replacement: "_"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	out := s.Apply(map[string]string{"KEY": "hello world!"})
	if out["KEY"] != "hello_world_" {
		t.Errorf("unexpected value: %q", out["KEY"])
	}
}

func TestApply_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := envsanitize.New(envsanitize.Options{
		Rules: []envsanitize.Rule{
			{Pattern: `[invalid`},
		},
	})
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestApply_DefaultOptions_AllEnabled(t *testing.T) {
	s, err := envsanitize.New(envsanitize.DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	out := s.Apply(map[string]string{
		"A": "  trim me  ",
		"B": "has\nnewline",
		"C": "null\x00byte",
	})
	if out["A"] != "trim me" {
		t.Errorf("A: got %q", out["A"])
	}
	if out["B"] != "has newline" {
		t.Errorf("B: got %q", out["B"])
	}
	if out["C"] != "nullbyte" {
		t.Errorf("C: got %q", out["C"])
	}
}

func TestApply_EmptySecrets(t *testing.T) {
	s, _ := envsanitize.New(envsanitize.DefaultOptions())
	out := s.Apply(map[string]string{})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestSanitized_SingleValue(t *testing.T) {
	s, _ := envsanitize.New(envsanitize.Options{TrimSpace: true, StripNewlines: true})
	got := s.Sanitized("  hello\nworld  ")
	if got != "hello world" {
		t.Errorf("got %q", got)
	}
}
