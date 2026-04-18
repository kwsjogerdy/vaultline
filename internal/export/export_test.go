package export

import (
	"bytes"
	"strings"
	"testing"
)

var sampleSecrets = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PASSWORD": "s3cr3t",
	"APP_NAME":    "my app",
}

func TestWrite_EnvFormat(t *testing.T) {
	var buf bytes.Buffer
	e := New(FormatEnv, &buf)
	if err := e.Write(sampleSecrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST=localhost in output, got:\n%s", out)
	}
	if !strings.Contains(out, `APP_NAME=`) {
		t.Errorf("expected APP_NAME in output, got:\n%s", out)
	}
}

func TestWrite_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	e := New(FormatJSON, &buf)
	if err := e.Write(sampleSecrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"DB_HOST"`) {
		t.Errorf("expected DB_HOST key in JSON output, got:\n%s", out)
	}
	if !strings.Contains(out, `"localhost"`) {
		t.Errorf("expected localhost value in JSON output, got:\n%s", out)
	}
}

func TestWrite_CSVFormat(t *testing.T) {
	var buf bytes.Buffer
	e := New(FormatCSV, &buf)
	if err := e.Write(sampleSecrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.HasPrefix(out, "key,value") {
		t.Errorf("expected CSV header, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in CSV output, got:\n%s", out)
	}
}

func TestWrite_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	e := New(Format("xml"), &buf)
	if err := e.Write(sampleSecrets); err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}

func TestWrite_SortedOutput(t *testing.T) {
	var buf bytes.Buffer
	e := New(FormatEnv, &buf)
	_ = e.Write(sampleSecrets)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "APP_NAME") {
		t.Errorf("expected sorted output, first line: %s", lines[0])
	}
}
