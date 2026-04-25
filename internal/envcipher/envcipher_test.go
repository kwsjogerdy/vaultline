package envcipher_test

import (
	"strings"
	"testing"

	"github.com/vaultline/vaultline/internal/envcipher"
)

func TestNew_EmptyPassphrase_ReturnsError(t *testing.T) {
	_, err := envcipher.New("", envcipher.AESGCM)
	if err == nil {
		t.Fatal("expected error for empty passphrase")
	}
}

func TestNew_UnsupportedAlgorithm_ReturnsError(t *testing.T) {
	_, err := envcipher.New("secret", "chacha20")
	if err == nil {
		t.Fatal("expected error for unsupported algorithm")
	}
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	c, err := envcipher.New("passphrase123", envcipher.AESGCM)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	plain := "super-secret-value"
	enc, err := c.Encrypt(plain)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if !strings.HasPrefix(enc, envcipher.Prefix) {
		t.Errorf("expected prefix %q, got %q", envcipher.Prefix, enc)
	}
	dec, err := c.Decrypt(enc)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if dec != plain {
		t.Errorf("got %q, want %q", dec, plain)
	}
}

func TestEncrypt_ProducesUniqueCiphertexts(t *testing.T) {
	c, _ := envcipher.New("pw", envcipher.AESGCM)
	a, _ := c.Encrypt("value")
	b, _ := c.Encrypt("value")
	if a == b {
		t.Error("expected unique ciphertexts due to random nonce")
	}
}

func TestDecrypt_NotEncrypted_ReturnsError(t *testing.T) {
	c, _ := envcipher.New("pw", envcipher.AESGCM)
	_, err := c.Decrypt("plaintext-no-prefix")
	if err == nil {
		t.Fatal("expected error for unencrypted value")
	}
}

func TestDecrypt_WrongKey_ReturnsError(t *testing.T) {
	c1, _ := envcipher.New("correct-key", envcipher.AESGCM)
	c2, _ := envcipher.New("wrong-key", envcipher.AESGCM)
	enc, _ := c1.Encrypt("secret")
	_, err := c2.Decrypt(enc)
	if err == nil {
		t.Fatal("expected decryption to fail with wrong key")
	}
}

func TestApplyEncrypt_TargetedKeys(t *testing.T) {
	c, _ := envcipher.New("pw", envcipher.AESGCM)
	secrets := map[string]string{
		"DB_PASS": "hunter2",
		"APP_ENV": "production",
	}
	out, err := c.ApplyEncrypt(secrets, []string{"DB_PASS"})
	if err != nil {
		t.Fatalf("ApplyEncrypt: %v", err)
	}
	if !strings.HasPrefix(out["DB_PASS"], envcipher.Prefix) {
		t.Error("DB_PASS should be encrypted")
	}
	if strings.HasPrefix(out["APP_ENV"], envcipher.Prefix) {
		t.Error("APP_ENV should not be encrypted")
	}
}

func TestApplyDecrypt_DecryptsAll(t *testing.T) {
	c, _ := envcipher.New("pw", envcipher.AESGCM)
	secrets := map[string]string{
		"DB_PASS": "hunter2",
		"API_KEY": "abc123",
	}
	encrypted, err := c.ApplyEncrypt(secrets, nil)
	if err != nil {
		t.Fatalf("ApplyEncrypt: %v", err)
	}
	decrypted, err := c.ApplyDecrypt(encrypted)
	if err != nil {
		t.Fatalf("ApplyDecrypt: %v", err)
	}
	for k, want := range secrets {
		if got := decrypted[k]; got != want {
			t.Errorf("key %q: got %q, want %q", k, got, want)
		}
	}
}
