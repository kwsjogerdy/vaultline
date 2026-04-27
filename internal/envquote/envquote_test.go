package envquote_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envquote"
)

func TestNew_UnknownStyle_ReturnsError(t *testing.T) {
	_, err := envquote.New("fancy")
	if err == nil {
		t.Fatal("expected error for unknown style")
	}
}

func TestNew_KnownStyles(t *testing.T) {
	for _, s := range []envquote.Style{envquote.StyleDouble, envquote.StyleSingle, envquote.StyleStrip} {
		_, err := envquote.New(s)
		if err != nil {
			t.Errorf("unexpected error for style %q: %v", s, err)
		}
	}
}

func TestApply_DoubleQuotes_WrapsValues(t *testing.T) {
	q, _ := envquote.New(envquote.StyleDouble)
	out := q.Apply(map[string]string{"KEY": "hello"})
	if got := out["KEY"]; got != `"hello"` {
		t.Errorf("got %q, want \"hello\"", got)
	}
}

func TestApply_SingleQuotes_WrapsValues(t *testing.T) {
	q, _ := envquote.New(envquote.StyleSingle)
	out := q.Apply(map[string]string{"KEY": "world"})
	if got := out["KEY"]; got != `'world'` {
		t.Errorf("got %q, want 'world'", got)
	}
}

func TestApply_Strip_RemovesDoubleQuotes(t *testing.T) {
	q, _ := envquote.New(envquote.StyleStrip)
	out := q.Apply(map[string]string{"KEY": `"value"`})
	if got := out["KEY"]; got != "value" {
		t.Errorf("got %q, want value", got)
	}
}

func TestApply_Strip_RemovesSingleQuotes(t *testing.T) {
	q, _ := envquote.New(envquote.StyleStrip)
	out := q.Apply(map[string]string{"KEY": `'value'`})
	if got := out["KEY"]; got != "value" {
		t.Errorf("got %q, want value", got)
	}
}

func TestApply_Double_ReplacesExistingSingleQuotes(t *testing.T) {
	q, _ := envquote.New(envquote.StyleDouble)
	out := q.Apply(map[string]string{"KEY": `'already'`})
	if got := out["KEY"]; got != `"already"` {
		t.Errorf("got %q, want \"already\"", got)
	}
}

func TestApply_PreservesAllKeys(t *testing.T) {
	q, _ := envquote.New(envquote.StyleStrip)
	secrets := map[string]string{"A": "1", "B": "2", "C": "3"}
	out := q.Apply(secrets)
	if len(out) != len(secrets) {
		t.Errorf("key count mismatch: got %d, want %d", len(out), len(secrets))
	}
}

func TestApply_EmptyValue_NoError(t *testing.T) {
	q, _ := envquote.New(envquote.StyleDouble)
	out := q.Apply(map[string]string{"EMPTY": ""})
	if got := out["EMPTY"]; got != `""` {
		t.Errorf("got %q, want empty double-quoted string", got)
	}
}
