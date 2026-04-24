// Package envencrypt provides per-key encryption for secret maps,
// allowing selective encryption of sensitive values at rest.
package envencrypt

import (
	"errors"
	"fmt"
	"strings"

	"github.com/vaultline/vaultline/internal/decrypt"
)

const encryptedPrefix = "enc:"

// Encrypter encrypts and decrypts individual secret values.
type Encrypter struct {
	passphrase string
	keys       map[string]struct{}
}

// New creates an Encrypter. keys is the set of secret keys whose values
// should be encrypted; an empty set means all keys are encrypted.
func New(passphrase string, keys []string) (*Encrypter, error) {
	if passphrase == "" {
		return nil, errors.New("envencrypt: passphrase must not be empty")
	}
	set := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		set[k] = struct{}{}
	}
	return &Encrypter{passphrase: passphrase, keys: set}, nil
}

// Encrypt returns a new map where targeted values are encrypted.
func (e *Encrypter) Encrypt(secrets map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if e.shouldProcess(k) {
			cipher, err := decrypt.Encrypt(v, e.passphrase)
			if err != nil {
				return nil, fmt.Errorf("envencrypt: encrypt key %q: %w", k, err)
			}
			out[k] = encryptedPrefix + cipher
		} else {
			out[k] = v
		}
	}
	return out, nil
}

// Decrypt returns a new map where enc:-prefixed values are decrypted.
func (e *Encrypter) Decrypt(secrets map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if strings.HasPrefix(v, encryptedPrefix) {
			plain, err := decrypt.Decrypt(strings.TrimPrefix(v, encryptedPrefix), e.passphrase)
			if err != nil {
				return nil, fmt.Errorf("envencrypt: decrypt key %q: %w", k, err)
			}
			out[k] = plain
		} else {
			out[k] = v
		}
	}
	return out, nil
}

// IsEncrypted reports whether the value for key k is encrypted.
func IsEncrypted(value string) bool {
	return strings.HasPrefix(value, encryptedPrefix)
}

func (e *Encrypter) shouldProcess(key string) bool {
	if len(e.keys) == 0 {
		return true
	}
	_, ok := e.keys[key]
	return ok
}
