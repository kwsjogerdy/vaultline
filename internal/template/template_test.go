package template

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env.tmpl")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	_ = f.Close()
	return f.Name()
}

func TestRenderBytes_BasicSubstitution(t *testing.T) {
	tmpl := []byte(`DB_HOST={{ index . "DB_HOST" }}`)
	out, err := RenderBytes(tmpl, map[string]string{"DB_HOST": "localhost"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "DB_HOST=localhost" {
		t.Errorf("got %q", out)
	}
}

func TestRenderBytes_DefaultFunc(t *testing.T) {
	tmpl := []byte(`PORT={{ default "8080" (index . "PORT") }}`)
	out, err := RenderBytes(tmpl, map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != `PORT=8080` {
		t.Errorf("got %q", out)
	}
}

func TestRenderBytes_RequiredFunc_Missing(t *testing.T) {
	tmpl := []byte(`SECRET={{ required "SECRET" (index . "SECRET") }}`)
	_, err := RenderBytes(tmpl, map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing required key")
	}
}

func TestRenderBytes_RequiredFunc_Present(t *testing.T) {
	tmpl := []byte(`SECRET={{ required "SECRET" (index . "SECRET") }}`)
	out, err := RenderBytes(tmpl, map[string]string{"SECRET": "abc123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "SECRET=abc123" {
		t.Errorf("got %q", out)
	}
}

func TestRender_FromFile(t *testing.T) {
	path := writeTemp(t, `APP_ENV={{ index . "APP_ENV" }}`)
	r := New(path)
	out, err := r.Render(map[string]string{"APP_ENV": "production"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "APP_ENV=production" {
		t.Errorf("got %q", out)
	}
}

func TestRender_MissingFile(t *testing.T) {
	r := New(filepath.Join(t.TempDir(), "nonexistent.tmpl"))
	_, err := r.Render(map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing template file")
	}
}
