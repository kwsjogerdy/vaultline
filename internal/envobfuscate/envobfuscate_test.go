package envobfuscate_test

import (
	"strings"
	"testing"

	"github.com/vaultline/vaultline/internal/envobfuscate"
)

func TestNew_EmptyKey_ReturnsError(t *testing.T) {
	_, err := envobfuscate.New("")
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestObfuscate_Deobfuscate_RoundTrip(t *testing.T) {
	o, _ := envobfuscate.New("secret-key")
	secrets := map[string]string{
		"DB_PASS": "hunter2",
		"API_KEY": "abc123xyz",
	}
	obf := o.Obfuscate(secrets)
	for k, v := range obf {
		if !envobfuscate.IsObfuscated(v) {
			t.Errorf("key %q value not marked as obfuscated", k)
		}
		if v == secrets[k] {
			t.Errorf("key %q value unchanged after obfuscation", k)
		}
	}
	restored, err := o.Deobfuscate(obf)
	if err != nil {
		t.Fatalf("deobfuscate: %v", err)
	}
	for k, want := range secrets {
		if got := restored[k]; got != want {
			t.Errorf("key %q: got %q, want %q", k, got, want)
		}
	}
}

func TestDeobfuscate_SkipsPlainValues(t *testing.T) {
	o, _ := envobfuscate.New("key")
	secrets := map[string]string{"PLAIN": "not-obfuscated"}
	out, err := o.Deobfuscate(secrets)
	if err != nil {
		t.Fatal(err)
	}
	if out["PLAIN"] != "not-obfuscated" {
		t.Errorf("plain value should pass through unchanged")
	}
}

func TestDeobfuscate_InvalidHex_ReturnsError(t *testing.T) {
	o, _ := envobfuscate.New("key")
	secrets := map[string]string{"BAD": "obf:ZZZZ"}
	_, err := o.Deobfuscate(secrets)
	if err == nil {
		t.Fatal("expected error for invalid hex payload")
	}
}

func TestIsObfuscated_True(t *testing.T) {
	if !envobfuscate.IsObfuscated("obf:deadbeef") {
		t.Error("expected true")
	}
}

func TestIsObfuscated_False(t *testing.T) {
	if envobfuscate.IsObfuscated("plaintext") {
		t.Error("expected false")
	}
}

func TestObfuscate_DifferentKeys_ProduceDifferentOutput(t *testing.T) {
	o1, _ := envobfuscate.New("key-one")
	o2, _ := envobfuscate.New("key-two")
	secrets := map[string]string{"X": "value"}
	out1 := o1.Obfuscate(secrets)["X"]
	out2 := o2.Obfuscate(secrets)["X"]
	if out1 == out2 {
		t.Error("different keys should produce different obfuscated output")
	}
	_ = strings.HasPrefix // suppress import warning
}
