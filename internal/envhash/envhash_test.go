package envhash_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envhash"
)

func TestCompute_StableAcrossKeyOrder(t *testing.T) {
	h := envhash.New()
	a := map[string]string{"FOO": "bar", "BAZ": "qux"}
	b := map[string]string{"BAZ": "qux", "FOO": "bar"}

	ra := h.Compute(a)
	rb := h.Compute(b)
	if ra.Fingerprint != rb.Fingerprint {
		t.Errorf("expected stable fingerprint, got %s vs %s", ra.Fingerprint, rb.Fingerprint)
	}
}

func TestCompute_DifferentValues_DifferentFingerprint(t *testing.T) {
	h := envhash.New()
	a := map[string]string{"FOO": "bar"}
	b := map[string]string{"FOO": "baz"}

	if h.Compute(a).Fingerprint == h.Compute(b).Fingerprint {
		t.Error("expected different fingerprints for different values")
	}
}

func TestCompute_EmptyMap(t *testing.T) {
	h := envhash.New()
	r := h.Compute(map[string]string{})
	if r.Fingerprint == "" {
		t.Error("expected non-empty fingerprint for empty map")
	}
	if r.KeyCount != 0 {
		t.Errorf("expected 0 keys, got %d", r.KeyCount)
	}
}

func TestCompute_KeysAreSorted(t *testing.T) {
	h := envhash.New()
	secrets := map[string]string{"Z_KEY": "1", "A_KEY": "2", "M_KEY": "3"}
	r := h.Compute(secrets)
	if r.Keys[0] != "A_KEY" || r.Keys[1] != "M_KEY" || r.Keys[2] != "Z_KEY" {
		t.Errorf("keys not sorted: %v", r.Keys)
	}
}

func TestEqual_SameMaps(t *testing.T) {
	h := envhash.New()
	a := map[string]string{"X": "1"}
	b := map[string]string{"X": "1"}
	if !h.Equal(a, b) {
		t.Error("expected Equal to return true for identical maps")
	}
}

func TestEqual_DifferentMaps(t *testing.T) {
	h := envhash.New()
	a := map[string]string{"X": "1"}
	b := map[string]string{"X": "2"}
	if h.Equal(a, b) {
		t.Error("expected Equal to return false for different maps")
	}
}
