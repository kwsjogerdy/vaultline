package decrypt_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/decrypt"
)

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	key := decrypt.DeriveKey("supersecretpassphrase")
	plaintext := "my_secret_value"

	encoded, err := decrypt.Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt error: %v", err)
	}

	result, err := decrypt.Decrypt(key, encoded)
	if err != nil {
		t.Fatalf("Decrypt error: %v", err)
	}
	if result != plaintext {
		t.Errorf("expected %q, got %q", plaintext, result)
	}
}

func TestEncrypt_ProducesUniqueCiphertexts(t *testing.T) {
	key := decrypt.DeriveKey("somekey")
	a, _ := decrypt.Encrypt(key, "value")
	b, _ := decrypt.Encrypt(key, "value")
	if a == b {
		t.Error("expected different ciphertexts due to random nonce")
	}
}

func TestDecrypt_InvalidBase64(t *testing.T) {
	key := decrypt.DeriveKey("somekey")
	_, err := decrypt.Decrypt(key, "!!!notbase64!!!")
	if err != decrypt.ErrInvalidCiphertext {
		t.Errorf("expected ErrInvalidCiphertext, got %v", err)
	}
}

func TestDecrypt_WrongKey(t *testing.T) {
	key1 := decrypt.DeriveKey("correctkey")
	key2 := decrypt.DeriveKey("wrongkeyxx")
	encoded, _ := decrypt.Encrypt(key1, "secret")
	_, err := decrypt.Decrypt(key2, encoded)
	if err != decrypt.ErrInvalidCiphertext {
		t.Errorf("expected ErrInvalidCiphertext, got %v", err)
	}
}

func TestDeriveKey_Length(t *testing.T) {
	key := decrypt.DeriveKey("short")
	if len(key) != 32 {
		t.Errorf("expected key length 32, got %d", len(key))
	}
}
