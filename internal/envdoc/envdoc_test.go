package envdoc_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/vaultline/vaultline/internal/envdoc"
)

func TestNew_EmptyKey_ReturnsError(t *testing.T) {
	_, err := envdoc.New(map[string]string{"": "desc|ex|required"})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestGet_Found(t *testing.T) {
	d, err := envdoc.New(map[string]string{
		"DB_URL": "Database URL|postgres://localhost/db|required",
	})
	if err != nil {
		t.Fatal(err)
	}
	e, ok := d.Get("DB_URL")
	if !ok {
		t.Fatal("expected entry to be found")
	}
	if e.Description != "Database URL" {
		t.Errorf("description = %q", e.Description)
	}
	if e.Example != "postgres://localhost/db" {
		t.Errorf("example = %q", e.Example)
	}
	if !e.Required {
		t.Error("expected Required=true")
	}
}

func TestGet_NotFound(t *testing.T) {
	d, _ := envdoc.New(map[string]string{})
	_, ok := d.Get("MISSING")
	if ok {
		t.Error("expected not found")
	}
}

func TestKeys_SortedOrder(t *testing.T) {
	d, _ := envdoc.New(map[string]string{
		"Z_KEY": "z|",
		"A_KEY": "a|",
		"M_KEY": "m|",
	})
	keys := d.Keys()
	if keys[0] != "A_KEY" || keys[1] != "M_KEY" || keys[2] != "Z_KEY" {
		t.Errorf("unexpected order: %v", keys)
	}
}

func TestWrite_TextFormat(t *testing.T) {
	d, _ := envdoc.New(map[string]string{
		"API_KEY": "API key|abc123|required",
	})
	var buf bytes.Buffer
	if err := d.Write(&buf, "text"); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "API_KEY") || !strings.Contains(out, "[required]") {
		t.Errorf("unexpected text output: %s", out)
	}
}

func TestWrite_MarkdownFormat(t *testing.T) {
	d, _ := envdoc.New(map[string]string{
		"SECRET": "A secret|s3cr3t|required",
	})
	var buf bytes.Buffer
	if err := d.Write(&buf, "markdown"); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "| Key |") || !strings.Contains(out, "SECRET") {
		t.Errorf("unexpected markdown output: %s", out)
	}
}

func TestWrite_UnsupportedFormat(t *testing.T) {
	d, _ := envdoc.New(map[string]string{})
	var buf bytes.Buffer
	if err := d.Write(&buf, "xml"); err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestNew_OptionalFields(t *testing.T) {
	d, err := envdoc.New(map[string]string{
		"SIMPLE": "just a description",
	})
	if err != nil {
		t.Fatal(err)
	}
	e, ok := d.Get("SIMPLE")
	if !ok {
		t.Fatal("expected entry")
	}
	if e.Required {
		t.Error("expected Required=false")
	}
	if e.Example != "" {
		t.Errorf("expected empty example, got %q", e.Example)
	}
}
