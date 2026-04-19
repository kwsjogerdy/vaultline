package envcast_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envcast"
)

func TestAsString_Found(t *testing.T) {
	c := envcast.New()
	secrets := map[string]string{"KEY": "hello"}
	v, err := c.AsString(secrets, "KEY")
	if err != nil || v != "hello" {
		t.Fatalf("expected hello, got %q err %v", v, err)
	}
}

func TestAsString_NotFound(t *testing.T) {
	c := envcast.New()
	_, err := c.AsString(map[string]string{}, "MISSING")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestAsInt_Valid(t *testing.T) {
	c := envcast.New()
	secrets := map[string]string{"PORT": "8080"}
	n, err := c.AsInt(secrets, "PORT")
	if err != nil || n != 8080 {
		t.Fatalf("expected 8080, got %d err %v", n, err)
	}
}

func TestAsInt_Invalid(t *testing.T) {
	c := envcast.New()
	secrets := map[string]string{"PORT": "not-a-number"}
	_, err := c.AsInt(secrets, "PORT")
	if err == nil {
		t.Fatal("expected error for invalid int")
	}
}

func TestAsBool_True(t *testing.T) {
	c := envcast.New()
	for _, v := range []string{"true", "1", "TRUE"} {
		b, err := c.AsBool(map[string]string{"F": v}, "F")
		if err != nil || !b {
			t.Fatalf("expected true for %q", v)
		}
	}
}

func TestAsBool_Invalid(t *testing.T) {
	c := envcast.New()
	_, err := c.AsBool(map[string]string{"F": "yes-please"}, "F")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAsFloat_Valid(t *testing.T) {
	c := envcast.New()
	f, err := c.AsFloat(map[string]string{"RATE": "3.14"}, "RATE")
	if err != nil || f != 3.14 {
		t.Fatalf("expected 3.14 got %f err %v", f, err)
	}
}

func TestCastAll_MixedTypes(t *testing.T) {
	c := envcast.New()
	secrets := map[string]string{"A": "42", "B": "true", "C": "hello"}
	types := map[string]string{"A": "int", "B": "bool"}
	out, errs := c.CastAll(secrets, types)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if out["A"] != int64(42) {
		t.Errorf("A: expected int64(42), got %v", out["A"])
	}
	if out["B"] != true {
		t.Errorf("B: expected true, got %v", out["B"])
	}
	if out["C"] != "hello" {
		t.Errorf("C: expected string, got %v", out["C"])
	}
}

func TestCastAll_CollectsErrors(t *testing.T) {
	c := envcast.New()
	secrets := map[string]string{"X": "bad"}
	_, errs := c.CastAll(secrets, map[string]string{"X": "int"})
	if len(errs) == 0 {
		t.Fatal("expected cast error")
	}
}
