// Package envchecksum provides checksum-based integrity verification for secret maps.
// It computes a stable SHA-256 checksum over a sorted key=value representation,
// allowing callers to detect unexpected mutations between operations.
package envchecksum

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// Checksum holds the hex-encoded SHA-256 digest of a secret map.
type Checksum struct {
	Digest string
	KeyCount int
}

// Verifier computes and verifies checksums for secret maps.
type Verifier struct{}

// New returns a new Verifier.
func New() *Verifier {
	return &Verifier{}
}

// Compute returns the Checksum for the given secrets map.
// Keys are sorted before hashing to ensure determinism.
func (v *Verifier) Compute(secrets map[string]string) Checksum {
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
	return Checksum{
		Digest:   hex.EncodeToString(sum[:]),
		KeyCount: len(keys),
	}
}

// Verify returns true when the checksum of secrets matches expected.
func (v *Verifier) Verify(secrets map[string]string, expected Checksum) bool {
	actual := v.Compute(secrets)
	return actual.Digest == expected.Digest
}

// String returns a human-readable representation of the checksum.
func (c Checksum) String() string {
	return fmt.Sprintf("sha256:%s (keys:%d)", c.Digest, c.KeyCount)
}
