package envtransform_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envtransform"
)

func TestApply_Upper(t *testing.T) {
	tr := envtransform.New()
	got, err := tr.Apply("upper", "hello")
	if err != nil {
		t.Fatal(err)
	}
	if got != "HELLO" {
		t.Errorf("expected HELLO, got %s", got)
	}
}

func TestApply_Lower(t *testing.T) {
	tr := envtransform.New()
	got, err := tr.Apply("lower", "WORLD")
	if err != nil {
		t.Fatal(err)
	}
	if got != "world" {
		t.Errorf("expected world, got %s", got)
	}
}

func TestApply_Trim(t *testing.T) {
	tr := envtransform.New()
	got, err := tr.Apply("trim", "  spaced  ")
	if err != nil {
		t.Fatal(err)
	}
	if got != "spaced" {
		t.Errorf("expected 'spaced', got %q", got)
	}
}

func TestApply_Unknown(t *testing.T) {
	tr := envtransform.New()
	_, err := tr.Apply("nonexistent", "value")
	if err == nil {
		t.Fatal("expected error for unknown transform")
	}
}

func TestApplyAll_ChainedTransforms(t *testing.T) {
	tr := envtransform.New()
	got, err := tr.ApplyAll([]string{"trim", "upper"}, "  hello  ")
	if err != nil {
		t.Fatal(err)
	}
	if got != "HELLO" {
		t.Errorf("expected HELLO, got %s", got)
	}
}

func TestApplyMap_AppliesPerKey(t *testing.T) {
	tr := envtransform.New()
	secrets := map[string]string{
		"DB_HOST": "  localhost  ",
		"APP_ENV": "production",
		"TOKEN":   "secret",
	}
	rules := map[string][]string{
		"DB_HOST": {"trim"},
		"APP_ENV": {"upper"},
	}
	out, err := tr.ApplyMap(secrets, rules)
	if err != nil {
		t.Fatal(err)
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected localhost, got %q", out["DB_HOST"])
	}
	if out["APP_ENV"] != "PRODUCTION" {
		t.Errorf("expected PRODUCTION, got %q", out["APP_ENV"])
	}
	if out["TOKEN"] != "secret" {
		t.Errorf("expected secret unchanged, got %q", out["TOKEN"])
	}
}

func TestRegister_CustomTransform(t *testing.T) {
	tr := envtransform.New()
	tr.Register("exclaim", func(v string) (string, error) { return v + "!", nil })
	got, err := tr.Apply("exclaim", "hello")
	if err != nil {
		t.Fatal(err)
	}
	if got != "hello!" {
		t.Errorf("expected hello!, got %s", got)
	}
}
