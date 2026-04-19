// Package envsign provides HMAC-based signing and verification for env secret maps.
package envsign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"sort"
)

// Signer signs and verifies secret maps.
type Signer struct {
	key []byte
}

// New creates a new Signer with the given secret key.
func New(key string) *Signer {
	return &Signer{key: []byte(key)}
}

// Sign returns an HMAC-SHA256 hex signature for the given secrets map.
// Keys are sorted before signing for determinism.
func (s *Signer) Sign(secrets map[string]string) (string, error) {
	payload, err := canonicalize(secrets)
	if err != nil {
		return "", err
	}
	mac := hmac.New(sha256.New, s.key)
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil)), nil
}

// Verify checks that the provided signature matches the secrets map.
func (s *Signer) Verify(secrets map[string]string, sig string) error {
	expected, err := s.Sign(secrets)
	if err != nil {
		return err
	}
	got, err := hex.DecodeString(sig)
	if err != nil {
		return errors.New("envsign: invalid signature encoding")
	}
	exp, _ := hex.DecodeString(expected)
	if !hmac.Equal(exp, got) {
		return errors.New("envsign: signature mismatch")
	}
	return nil
}

// canonicalize produces a deterministic JSON encoding of the secrets map.
func canonicalize(secrets map[string]string) ([]byte, error) {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	ordered := make([][2]string, 0, len(keys))
	for _, k := range keys {
		ordered = append(ordered, [2]string{k, secrets[k]})
	}
	return json.Marshal(ordered)
}
