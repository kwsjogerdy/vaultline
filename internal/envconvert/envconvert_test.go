package envconvert

import (
	"encoding/base64"
	"testing"
)

func TestNew_UnknownFormat(t *testing.T) {
	_, err := New("hex")
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestNew_KnownFormat(t *testing.T) {
	c, err := New(FormatBase64)
	if err != nil || c == nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApply_Base64(t *testing.T) {
	c, _ := New(FormatBase64)
	_, out := c.Apply(map[string]string{"KEY": "hello"})
	want := base64.StdEncoding.EncodeToString([]byte("hello"))
	if out["KEY"] != want {
		t.Errorf("got %q want %q", out["KEY"], want)
	}
}

func TestApply_URLEncode(t *testing.T) {
	c, _ := New(FormatURLEncode)
	_, out := c.Apply(map[string]string{"K": "hello world"})
	if out["K"] != "hello+world" {
		t.Errorf("got %q", out["K"])
	}
}

func TestApply_Uppercase(t *testing.T) {
	c, _ := New(FormatUppercase)
	_, out := c.Apply(map[string]string{"K": "secret"})
	if out["K"] != "SECRET" {
		t.Errorf("got %q", out["K"])
	}
}

func TestApply_Lowercase(t *testing.T) {
	c, _ := New(FormatLowercase)
	_, out := c.Apply(map[string]string{"K": "VALUE"})
	if out["K"] != "value" {
		t.Errorf("got %q", out["K"])
	}
}

func TestApply_TrimSpace(t *testing.T) {
	c, _ := New(FormatTrimSpace)
	_, out := c.Apply(map[string]string{"K": "  trimmed  "})
	if out["K"] != "trimmed" {
		t.Errorf("got %q", out["K"])
	}
}

func TestApply_ReturnsResults(t *testing.T) {
	c, _ := New(FormatUppercase)
	results, _ := c.Apply(map[string]string{"A": "x"})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Original != "x" || results[0].Converted != "X" {
		t.Errorf("unexpected result: %+v", results[0])
	}
}
