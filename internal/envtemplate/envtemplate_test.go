package envtemplate_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envtemplate"
)

func TestRender_SimpleSubstitution(t *testing.T) {
	r := envtemplate.New(map[string]string{"HOST": "localhost", "PORT": "5432"})
	out, err := r.Render("postgres://{{index . \"HOST\"}}:{{index . \"PORT\"}}/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "postgres://localhost:5432/db" {
		t.Errorf("got %q", out)
	}
}

func TestRender_SecretFunc(t *testing.T) {
	r := envtemplate.New(map[string]string{"API_KEY": "abc123"})
	out, err := r.Render(`{{ secret "API_KEY" }}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "abc123" {
		t.Errorf("got %q", out)
	}
}

func TestRender_SecretFunc_Missing(t *testing.T) {
	r := envtemplate.New(map[string]string{})
	_, err := r.Render(`{{ secret "MISSING" }}`)
	if err == nil {
		t.Fatal("expected error for missing secret")
	}
}

func TestRender_MissingKey_ReturnsError(t *testing.T) {
	r := envtemplate.New(map[string]string{"A": "1"})
	_, err := r.Render(`{{index . "NOPE"}}`)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRender_UpperHelper(t *testing.T) {
	r := envtemplate.New(map[string]string{"ENV": "production"})
	out, err := r.Render(`{{ secret "ENV" | upper }}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "PRODUCTION" {
		t.Errorf("got %q", out)
	}
}

func TestRenderMap_ExpandsValues(t *testing.T) {
	r := envtemplate.New(map[string]string{
		"BASE_URL": "https://example.com",
		"CALLBACK": `{{ secret "BASE_URL" }}/callback`,
	})
	out, err := r.RenderMap()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["CALLBACK"] != "https://example.com/callback" {
		t.Errorf("got %q", out["CALLBACK"])
	}
}

func TestRenderMap_ErrorPropagates(t *testing.T) {
	r := envtemplate.New(map[string]string{
		"BAD": `{{ secret "UNDEFINED" }}`,
	})
	_, err := r.RenderMap()
	if err == nil {
		t.Fatal("expected error")
	}
}
