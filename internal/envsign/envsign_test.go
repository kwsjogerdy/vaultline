package envsign_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envsign"
)

func TestSign_Deterministic(t *testing.T) {
	s := envsign.New("supersecret")
	secrets := map[string]string{"B": "2", "A": "1", "C": "3"}
	sig1, err := s.Sign(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sig2, err := s.Sign(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sig1 != sig2 {
		t.Errorf("expected deterministic signatures, got %q and %q", sig1, sig2)
	}
}

func TestVerify_ValidSignature(t *testing.T) {
	s := envsign.New("mykey")
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	sig, err := s.Sign(secrets)
	if err != nil {
		t.Fatalf("sign error: %v", err)
	}
	if err := s.Verify(secrets, sig); err != nil {
		t.Errorf("expected valid signature, got error: %v", err)
	}
}

func TestVerify_TamperedSecrets(t *testing.T) {
	s := envsign.New("mykey")
	secrets := map[string]string{"TOKEN": "abc"}
	sig, _ := s.Sign(secrets)
	tampered := map[string]string{"TOKEN": "xyz"}
	if err := s.Verify(tampered, sig); err == nil {
		t.Error("expected signature mismatch error")
	}
}

func TestVerify_WrongKey(t *testing.T) {
	s1 := envsign.New("key1")
	s2 := envsign.New("key2")
	secrets := map[string]string{"X": "1"}
	sig, _ := s1.Sign(secrets)
	if err := s2.Verify(secrets, sig); err == nil {
		t.Error("expected mismatch with different key")
	}
}

func TestVerify_InvalidHex(t *testing.T) {
	s := envsign.New("key")
	if err := s.Verify(map[string]string{}, "not-hex!"); err == nil {
		t.Error("expected error for invalid hex signature")
	}
}

func TestSign_EmptySecrets(t *testing.T) {
	s := envsign.New("key")
	sig, err := s.Sign(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sig == "" {
		t.Error("expected non-empty signature for empty map")
	}
}
