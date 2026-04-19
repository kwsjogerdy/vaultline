package envflatten

import (
	"testing"
)

func TestFlatten_SimpleMap(t *testing.T) {
	f := New("__")
	out, err := f.Flatten(map[string]any{"db": "postgres", "port": "5432"})
	if err != nil {
		t.Fatal(err)
	}
	if out["DB"] != "postgres" || out["PORT"] != "5432" {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestFlatten_Nested(t *testing.T) {
	f := New("__")
	out, err := f.Flatten(map[string]any{
		"database": map[string]any{
			"host": "localhost",
			"port": "5432",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if out["DATABASE__HOST"] != "localhost" {
		t.Errorf("expected DATABASE__HOST=localhost, got %v", out)
	}
	if out["DATABASE__PORT"] != "5432" {
		t.Errorf("expected DATABASE__PORT=5432, got %v", out)
	}
}

func TestFlatten_WithPrefix(t *testing.T) {
	f := New("__").WithPrefix("app")
	out, err := f.Flatten(map[string]any{"name": "vaultline"})
	if err != nil {
		t.Fatal(err)
	}
	if out["APP__NAME"] != "vaultline" {
		t.Errorf("expected APP__NAME=vaultline, got %v", out)
	}
}

func TestFlatten_NumericValue(t *testing.T) {
	f := New("__")
	out, err := f.Flatten(map[string]any{"timeout": 30})
	if err != nil {
		t.Fatal(err)
	}
	if out["TIMEOUT"] != "30" {
		t.Errorf("expected TIMEOUT=30, got %v", out)
	}
}

func TestFlatten_NilValue(t *testing.T) {
	f := New("__")
	out, err := f.Flatten(map[string]any{"empty": nil})
	if err != nil {
		t.Fatal(err)
	}
	if out["EMPTY"] != "" {
		t.Errorf("expected EMPTY='', got %v", out)
	}
}

func TestFlatten_UnsupportedType(t *testing.T) {
	f := New("__")
	_, err := f.Flatten(map[string]any{"bad": []string{"a", "b"}})
	if err == nil {
		t.Error("expected error for unsupported type")
	}
}

func TestFlatten_CustomSeparator(t *testing.T) {
	f := New(".")
	out, err := f.Flatten(map[string]any{"db": map[string]any{"host": "rds"}})
	if err != nil {
		t.Fatal(err)
	}
	if out["DB.HOST"] != "rds" {
		t.Errorf("expected DB.HOST=rds, got %v", out)
	}
}
