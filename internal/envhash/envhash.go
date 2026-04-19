// Package envhash computes a stable hash fingerprint of a secrets map.
// Useful for detecting changes between syncs without storing raw values.
package envhash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// Result holds the fingerprint of a secrets map.
type Result struct {
	Fingerprint string
	KeyCount    int
	Keys        []string
}

// Hasher computes fingerprints for secret maps.
type Hasher struct{}

// New returns a new Hasher.
func New() *Hasher {
	return &Hasher{}
}

// Compute returns a SHA-256 fingerprint of the provided secrets map.
// Keys are sorted before hashing to ensure stability.
func (h *Hasher) Compute(secrets map[string]string) Result {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%s\n", k, secrets[k])
	}

	sum := sha256.Sum256([]byte(sb.String()))
	return Result{
		Fingerprint: hex.EncodeToString(sum[:]),
		KeyCount:    len(keys),
		Keys:        keys,
	}
}

// Equal returns true if two secret maps produce the same fingerprint.
func (h *Hasher) Equal(a, b map[string]string) bool {
	return h.Compute(a).Fingerprint == h.Compute(b).Fingerprint
}
