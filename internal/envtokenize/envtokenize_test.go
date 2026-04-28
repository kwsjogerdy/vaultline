package envtokenize

import (
	"testing"
)

func TestNew_EmptyDelimiter_ReturnsError(t *testing.T) {
	_, err := New("", false)
	if err == nil {
		t.Fatal("expected error for empty delimiter")
	}
}

func TestApply_SplitsByDelimiter(t *testing.T) {
	tok, _ := New(",", false)
	secrets := map[string]string{
		"HOSTS": "a,b,c",
		"SINGLE": "only",
	}
	results := tok.Apply(secrets)

	if len(results["HOSTS"].Tokens) != 3 {
		t.Fatalf("expected 3 tokens, got %d", len(results["HOSTS"].Tokens))
	}
	if results["HOSTS"].Tokens[1] != "b" {
		t.Errorf("expected 'b', got %q", results["HOSTS"].Tokens[1])
	}
	if len(results["SINGLE"].Tokens) != 1 {
		t.Fatalf("expected 1 token for SINGLE, got %d", len(results["SINGLE"].Tokens))
	}
}

func TestApply_TrimEnabled(t *testing.T) {
	tok, _ := New(",", true)
	secrets := map[string]string{"LIST": " foo , bar , baz "}
	results := tok.Apply(secrets)
	tokens := results["LIST"].Tokens
	expected := []string{"foo", "bar", "baz"}
	for i, e := range expected {
		if tokens[i] != e {
			t.Errorf("token[%d]: expected %q, got %q", i, e, tokens[i])
		}
	}
}

func TestApply_RawPreserved(t *testing.T) {
	tok, _ := New(":", false)
	secrets := map[string]string{"CONN": "host:5432"}
	results := tok.Apply(secrets)
	if results["CONN"].Raw != "host:5432" {
		t.Errorf("expected raw value preserved, got %q", results["CONN"].Raw)
	}
}

func TestGet_ValidIndex(t *testing.T) {
	tok, _ := New(":", false)
	secrets := map[string]string{"ADDR": "localhost:8080"}
	v, ok := tok.Get(secrets, "ADDR", 1)
	if !ok {
		t.Fatal("expected ok=true")
	}
	if v != "8080" {
		t.Errorf("expected '8080', got %q", v)
	}
}

func TestGet_OutOfRange(t *testing.T) {
	tok, _ := New(":", false)
	secrets := map[string]string{"ADDR": "localhost:8080"}
	_, ok := tok.Get(secrets, "ADDR", 5)
	if ok {
		t.Error("expected ok=false for out-of-range index")
	}
}

func TestGet_MissingKey(t *testing.T) {
	tok, _ := New(",", false)
	secrets := map[string]string{}
	_, ok := tok.Get(secrets, "MISSING", 0)
	if ok {
		t.Error("expected ok=false for missing key")
	}
}

func TestApply_EmptySecrets(t *testing.T) {
	tok, _ := New(",", false)
	results := tok.Apply(map[string]string{})
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d entries", len(results))
	}
}
