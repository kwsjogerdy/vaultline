package envencrypt_test

import (
	"strings"
	"testing"

	"github.com/vaultline/vaultline/internal/envencrypt"
)

func TestNew_EmptyPassphrase_ReturnsError(t *testing.T) {
	_, err := envencrypt.New("", nil)
	if err == nil {
		t.Fatal("expected error for empty passphrase")
	}
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	e, err := envencrypt.New("s3cret", nil)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	secrets := map[string]string{
		"DB_PASSWORD": "hunter2",
		"API_KEY":     "abc123",
	}
	enc, err := e.Encrypt(secrets)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	for k, v := range enc {
		if !envencrypt.IsEncrypted(v) {
			t.Errorf("key %q value not marked encrypted: %q", k, v)
		}
	}
	dec, err := e.Decrypt(enc)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	for k, want := range secrets {
		if got := dec[k]; got != want {
			t.Errorf("key %q: got %q want %q", k, got, want)
		}
	}
}

func TestEncrypt_OnlyTargetedKeys(t *testing.T) {
	e, _ := envencrypt.New("pass", []string{"SECRET"})
	secrets := map[string]string{
		"SECRET": "topsecret",
		"PUBLIC": "visible",
	}
	enc, err := e.Encrypt(secrets)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if !envencrypt.IsEncrypted(enc["SECRET"]) {
		t.Error("SECRET should be encrypted")
	}
	if envencrypt.IsEncrypted(enc["PUBLIC"]) {
		t.Error("PUBLIC should not be encrypted")
	}
}

func TestDecrypt_SkipsPlainValues(t *testing.T) {
	e, _ := envencrypt.New("pass", nil)
	secrets := map[string]string{"PLAIN": "hello"}
	dec, err := e.Decrypt(secrets)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if dec["PLAIN"] != "hello" {
		t.Errorf("expected plain value unchanged, got %q", dec["PLAIN"])
	}
}

func TestDecrypt_WrongPassphrase_ReturnsError(t *testing.T) {
	enc, _ := envencrypt.New("correct", nil)
	wrong, _ := envencrypt.New("wrong", nil)
	encrypted, err := enc.Encrypt(map[string]string{"K": "v"})
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	_, err = wrong.Decrypt(encrypted)
	if err == nil {
		t.Fatal("expected error when decrypting with wrong passphrase")
	}
}

func TestIsEncrypted(t *testing.T) {
	if !envencrypt.IsEncrypted("enc:abc") {
		t.Error("expected true for enc: prefix")
	}
	if envencrypt.IsEncrypted("plaintext") {
		t.Error("expected false for plain value")
	}
	if envencrypt.IsEncrypted("") {
		t.Error("expected false for empty string")
	}
}

func TestEncrypt_EmptySecrets(t *testing.T) {
	e, _ := envencrypt.New("pass", nil)
	out, err := e.Encrypt(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestEncrypt_ProducesEncPrefix(t *testing.T) {
	e, _ := envencrypt.New("pass", nil)
	out, _ := e.Encrypt(map[string]string{"X": "val"})
	if !strings.HasPrefix(out["X"], "enc:") {
		t.Errorf("expected enc: prefix, got %q", out["X"])
	}
}
