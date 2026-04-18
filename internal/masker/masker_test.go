package masker_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/masker"
)

func TestMask_FullRedact(t *testing.T) {
	m := masker.New()
	if got := m.Mask("supersecret"); got != "****" {
		t.Fatalf("expected ****, got %s", got)
	}
}

func TestMask_EmptyValue(t *testing.T) {
	m := masker.New()
	if got := m.Mask(""); got != "****" {
		t.Fatalf("expected **** for empty, got %s", got)
	}
}

func TestMask_RevealPrefix(t *testing.T) {
	m := masker.NewWithReveal(3)
	got := m.Mask("abcdef")
	if got != "abc****" {
		t.Fatalf("expected abc****, got %s", got)
	}
}

func TestMask_RevealExceedsLength(t *testing.T) {
	m := masker.NewWithReveal(20)
	if got := m.Mask("short"); got != "****" {
		t.Fatalf("expected **** when reveal >= len, got %s", got)
	}
}

func TestMaskMap_RedactsAllValues(t *testing.T) {
	m := masker.New()
	input := map[string]string{"KEY1": "val1", "KEY2": "val2"}
	out := m.MaskMap(input)
	for k, v := range out {
		if v != "****" {
			t.Errorf("key %s: expected ****, got %s", k, v)
		}
	}
	if len(out) != len(input) {
		t.Errorf("map length mismatch")
	}
}

func TestMaskLine_ReplacesSecrets(t *testing.T) {
	m := masker.New()
	secrets := map[string]string{"TOKEN": "abc123"}
	line := "export TOKEN=abc123"
	got := m.MaskLine(line, secrets)
	expected := "export TOKEN=****"
	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestMaskLine_NoMatch(t *testing.T) {
	m := masker.New()
	secrets := map[string]string{"TOKEN": "xyz"}
	line := "nothing here"
	if got := m.MaskLine(line, secrets); got != "nothing here" {
		t.Fatalf("unexpected change: %s", got)
	}
}
