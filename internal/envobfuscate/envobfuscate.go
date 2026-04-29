// Package envobfuscate provides reversible obfuscation of secret values
// using a deterministic substitution cipher suitable for logging and display.
package envobfuscate

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

const obfuscatedPrefix = "obf:"

// Obfuscator encodes and decodes secret values with a shared key.
type Obfuscator struct {
	key []byte
}

// New creates an Obfuscator with the given key.
func New(key string) (*Obfuscator, error) {
	if strings.TrimSpace(key) == "" {
		return nil, errors.New("envobfuscate: key must not be empty")
	}
	return &Obfuscator{key: []byte(key)}, nil
}

// Obfuscate encodes all values in secrets, returning a new map.
func (o *Obfuscator) Obfuscate(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = obfuscatedPrefix + o.encode(v)
	}
	return out
}

// Deobfuscate decodes all obfuscated values in secrets, returning a new map.
func (o *Obfuscator) Deobfuscate(secrets map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if !strings.HasPrefix(v, obfuscatedPrefix) {
			out[k] = v
			continue
		}
		decoded, err := o.decode(strings.TrimPrefix(v, obfuscatedPrefix))
		if err != nil {
			return nil, fmt.Errorf("envobfuscate: decode key %q: %w", k, err)
		}
		out[k] = decoded
	}
	return out, nil
}

// IsObfuscated reports whether the value was produced by Obfuscate.
func IsObfuscated(value string) bool {
	return strings.HasPrefix(value, obfuscatedPrefix)
}

// encode XORs each byte of value with the HMAC-derived keystream.
func (o *Obfuscator) encode(value string) string {
	stream := o.keystream(len(value))
	b := []byte(value)
	for i := range b {
		b[i] ^= stream[i]
	}
	return hex.EncodeToString(b)
}

func (o *Obfuscator) decode(encoded string) (string, error) {
	b, err := hex.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	stream := o.keystream(len(b))
	for i := range b {
		b[i] ^= stream[i]
	}
	return string(b), nil
}

// keystream generates n bytes derived from the HMAC-SHA256 of the key.
func (o *Obfuscator) keystream(n int) []byte {
	stream := make([]byte, 0, n)
	counter := byte(0)
	for len(stream) < n {
		mac := hmac.New(sha256.New, o.key)
		mac.Write([]byte{counter})
		stream = append(stream, mac.Sum(nil)...)
		counter++
	}
	return stream[:n]
}
