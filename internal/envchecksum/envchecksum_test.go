package envchecksum_test

import (
	"strings"
	"testing"

	"github.com/vaultline/vaultline/internal/envchecksum"
)

func TestCompute_StableAcrossKeyOrder(t *testing.T) {
	v := envchecksum.New()
	a := map[string]string{"FOO": "bar", "BAZ": "qux"}
	b := map[string]string{"BAZ": "qux", "FOO": "bar"}

	ca := v.Compute(a)
	cb := v.Compute(b)

	if ca.Digest != cb.Digest {
		t.Errorf("expected same digest for same content, got %s vs %s", ca.Digest, cb.Digest)
	}
}

func TestCompute_DifferentValues_DifferentDigest(t *testing.T) {
	v := envchecksum.New()
	a := map[string]string{"FOO": "bar"}
	b := map[string]string{"FOO": "baz"}

	if v.Compute(a).Digest == v.Compute(b).Digest {
		t.Error("expected different digests for different values")
	}
}

func TestCompute_EmptyMap(t *testing.T) {
	v := envchecksum.New()
	c := v.Compute(map[string]string{})

	if c.Digest == "" {
		t.Error("expected non-empty digest for empty map")
	}
	if c.KeyCount != 0 {
		t.Errorf("expected KeyCount 0, got %d", c.KeyCount)
	}
}

func TestCompute_KeyCountIsCorrect(t *testing.T) {
	v := envchecksum.New()
	secrets := map[string]string{"A": "1", "B": "2", "C": "3"}
	c := v.Compute(secrets)

	if c.KeyCount != 3 {
		t.Errorf("expected KeyCount 3, got %d", c.KeyCount)
	}
}

func TestVerify_MatchingChecksum(t *testing.T) {
	v := envchecksum.New()
	secrets := map[string]string{"TOKEN": "abc123"}
	c := v.Compute(secrets)

	if !v.Verify(secrets, c) {
		t.Error("expected Verify to return true for matching checksum")
	}
}

func TestVerify_MutatedSecrets(t *testing.T) {
	v := envchecksum.New()
	original := map[string]string{"TOKEN": "abc123"}
	c := v.Compute(original)

	mutated := map[string]string{"TOKEN": "xyz999"}
	if v.Verify(mutated, c) {
		t.Error("expected Verify to return false for mutated secrets")
	}
}

func TestChecksum_String_ContainsDigest(t *testing.T) {
	v := envchecksum.New()
	c := v.Compute(map[string]string{"K": "V"})
	s := c.String()

	if !strings.HasPrefix(s, "sha256:") {
		t.Errorf("expected String() to start with sha256:, got %q", s)
	}
	if !strings.Contains(s, "keys:1") {
		t.Errorf("expected String() to contain key count, got %q", s)
	}
}
