package envreport_test

import (
	"strings"
	"testing"

	"github.com/vaultline/vaultline/internal/envreport"
)

func TestBuild_CountsTotal(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2", "C": "3"}
	rep := envreport.Build(secrets)
	if rep.Total != 3 {
		t.Fatalf("expected 3, got %d", rep.Total)
	}
}

func TestBuild_SortsKeys(t *testing.T) {
	secrets := map[string]string{"Z_KEY": "1", "A_KEY": "2", "M_KEY": "3"}
	rep := envreport.Build(secrets)
	if rep.Keys[0] != "A_KEY" {
		t.Fatalf("expected A_KEY first, got %s", rep.Keys[0])
	}
}

func TestBuild_GroupsByPrefix(t *testing.T) {
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_KEY": "secret",
	}
	rep := envreport.Build(secrets)
	if rep.BySuffix["DB"] != 2 {
		t.Fatalf("expected DB count 2, got %d", rep.BySuffix["DB"])
	}
	if rep.BySuffix["APP"] != 1 {
		t.Fatalf("expected APP count 1, got %d", rep.BySuffix["APP"])
	}
}

func TestWrite_TextFormat(t *testing.T) {
	var buf strings.Builder
	r := envreport.New(&buf, envreport.FormatText)
	rep := envreport.Build(map[string]string{"DB_HOST": "localhost", "APP_KEY": "x"})
	if err := r.Write(rep); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "Total secrets: 2") {
		t.Errorf("missing total line: %s", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("missing key DB_HOST: %s", out)
	}
}

func TestWrite_MarkdownFormat(t *testing.T) {
	var buf strings.Builder
	r := envreport.New(&buf, envreport.FormatMarkdown)
	rep := envreport.Build(map[string]string{"DB_HOST": "localhost"})
	if err := r.Write(rep); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "## Secret Report") {
		t.Errorf("missing markdown header: %s", out)
	}
	if !strings.Contains(out, "**Total:**") {
		t.Errorf("missing bold total: %s", out)
	}
}

func TestBuild_EmptySecrets(t *testing.T) {
	rep := envreport.Build(map[string]string{})
	if rep.Total != 0 {
		t.Fatalf("expected 0, got %d", rep.Total)
	}
}
