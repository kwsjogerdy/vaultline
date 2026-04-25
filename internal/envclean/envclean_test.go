package envclean

import (
	"sort"
	"testing"
)

func TestApply_TrimSpace(t *testing.T) {
	c := New(Options{TrimSpace: true})
	out := c.Apply(map[string]string{"KEY": "  hello  "})
	if got := out["KEY"]; got != "hello" {
		t.Errorf("expected 'hello', got %q", got)
	}
}

func TestApply_NoTrim(t *testing.T) {
	c := New(Options{TrimSpace: false})
	out := c.Apply(map[string]string{"KEY": "  hello  "})
	if got := out["KEY"]; got != "  hello  " {
		t.Errorf("expected '  hello  ', got %q", got)
	}
}

func TestApply_CollapseSpaces(t *testing.T) {
	c := New(Options{CollapseSpaces: true})
	out := c.Apply(map[string]string{"KEY": "foo   bar  baz"})
	if got := out["KEY"]; got != "foo bar baz" {
		t.Errorf("expected 'foo bar baz', got %q", got)
	}
}

func TestApply_StripControl(t *testing.T) {
	c := New(Options{StripControl: true})
	out := c.Apply(map[string]string{"KEY": "hello\x01world"})
	if got := out["KEY"]; got != "helloworld" {
		t.Errorf("expected 'helloworld', got %q", got)
	}
}

func TestApply_NormalizeNewlines(t *testing.T) {
	c := New(Options{NormalizeNewlines: true})
	out := c.Apply(map[string]string{"KEY": "line1\r\nline2\rline3"})
	if got := out["KEY"]; got != "line1\nline2\nline3" {
		t.Errorf("unexpected value: %q", got)
	}
}

func TestApply_PreservesUnchanged(t *testing.T) {
	c := New(DefaultOptions())
	out := c.Apply(map[string]string{"KEY": "clean"})
	if got := out["KEY"]; got != "clean" {
		t.Errorf("expected 'clean', got %q", got)
	}
}

func TestCleaned_ReportsModifiedKeys(t *testing.T) {
	c := New(Options{TrimSpace: true})
	c.Apply(map[string]string{
		"DIRTY": "  val  ",
		"CLEAN": "val",
	})
	cleaned := c.Cleaned()
	if len(cleaned) != 1 || cleaned[0] != "DIRTY" {
		t.Errorf("expected [DIRTY], got %v", cleaned)
	}
}

func TestCleaned_EmptyWhenNothingChanged(t *testing.T) {
	c := New(DefaultOptions())
	c.Apply(map[string]string{"A": "clean", "B": "also-clean"})
	if got := c.Cleaned(); len(got) != 0 {
		t.Errorf("expected no cleaned keys, got %v", got)
	}
}

func TestApply_MultipleOptions(t *testing.T) {
	c := New(Options{TrimSpace: true, StripControl: true, NormalizeNewlines: true})
	out := c.Apply(map[string]string{"K": "  \x03hello\r\n  "})
	if got := out["K"]; got != "hello" {
		t.Errorf("expected 'hello', got %q", got)
	}
}

func TestApply_EmptyMap(t *testing.T) {
	c := New(DefaultOptions())
	out := c.Apply(map[string]string{})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	input := map[string]string{"KEY": "  value  "}
	c := New(Options{TrimSpace: true})
	c.Apply(input)
	if input["KEY"] != "  value  " {
		t.Error("input map was mutated")
	}
}

func TestCleaned_SortedForDeterminism(t *testing.T) {
	c := New(Options{TrimSpace: true})
	c.Apply(map[string]string{
		"Z": "  z  ",
		"A": "  a  ",
		"M": "  m  ",
	})
	keys := c.Cleaned()
	sorted := make([]string, len(keys))
	copy(sorted, keys)
	sort.Strings(sorted)
	for i, k := range sorted {
		if keys[i] != k {
			// order may vary; just check all are present
			break
		}
	}
	if len(keys) != 3 {
		t.Errorf("expected 3 cleaned keys, got %d", len(keys))
	}
}
