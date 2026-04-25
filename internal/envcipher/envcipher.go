// Package envcipher provides field-level encryption for individual secret values
// using a named cipher strategy (AES-GCM or ChaCha20-Poly1305).
package envcipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Algorithm identifies the cipher to use.
type Algorithm string

const (
	AESGCM    Algorithm = "aes-gcm"
	Prefix              = "enc:"
)

// Cipher encrypts and decrypts individual secret values.
type Cipher struct {
	key []byte
	alg Algorithm
}

// New creates a Cipher with the given passphrase and algorithm.
// Only "aes-gcm" is currently supported.
func New(passphrase string, alg Algorithm) (*Cipher, error) {
	if passphrase == "" {
		return nil, errors.New("envcipher: passphrase must not be empty")
	}
	if alg != AESGCM {
		return nil, fmt.Errorf("envcipher: unsupported algorithm %q", alg)
	}
	h := sha256.Sum256([]byte(passphrase))
	return &Cipher{key: h[:], alg: alg}, nil
}

// Encrypt encrypts plaintext and returns a prefixed base64 ciphertext.
func (c *Cipher) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	sealed := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return Prefix + base64.StdEncoding.EncodeToString(sealed), nil
}

// Decrypt decrypts a value produced by Encrypt.
func (c *Cipher) Decrypt(ciphertext string) (string, error) {
	if !strings.HasPrefix(ciphertext, Prefix) {
		return "", fmt.Errorf("envcipher: value is not encrypted")
	}
	raw, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(ciphertext, Prefix))
	if err != nil {
		return "", fmt.Errorf("envcipher: invalid base64: %w", err)
	}
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	ns := gcm.NonceSize()
	if len(raw) < ns {
		return "", errors.New("envcipher: ciphertext too short")
	}
	plain, err := gcm.Open(nil, raw[:ns], raw[ns:], nil)
	if err != nil {
		return "", fmt.Errorf("envcipher: decryption failed: %w", err)
	}
	return string(plain), nil
}

// ApplyEncrypt returns a copy of secrets with the given keys encrypted.
func (c *Cipher) ApplyEncrypt(secrets map[string]string, keys []string) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	target := toSet(keys)
	for k, v := range secrets {
		if len(target) == 0 || target[k] {
			enc, err := c.Encrypt(v)
			if err != nil {
				return nil, fmt.Errorf("envcipher: key %q: %w", k, err)
			}
			out[k] = enc
		} else {
			out[k] = v
		}
	}
	return out, nil
}

// ApplyDecrypt returns a copy of secrets with all encrypted values decrypted.
func (c *Cipher) ApplyDecrypt(secrets map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if strings.HasPrefix(v, Prefix) {
			dec, err := c.Decrypt(v)
			if err != nil {
				return nil, fmt.Errorf("envcipher: key %q: %w", k, err)
			}
			out[k] = dec
		} else {
			out[k] = v
		}
	}
	return out, nil
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
