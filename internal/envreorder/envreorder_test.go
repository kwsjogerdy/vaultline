package envreorder_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envreorder"
)

func TestNew_EmptyOrderAndNoFallback_ReturnsError(t *testing.T) {
	_, err := envreorder.New(nil, false)
	if err == nil {
		t.Fatal("expected error for empty order and no fallback")
	}
}

func TestNew_EmptyOrderWithFallback_OK(t *testing.T) {
	_, err := envreorder.New(nil, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApply_ExplicitOrder(t *testing.T) {
	r, _ := envreorder.New([]string{"C", "A", "B"}, false)
	secrets := map[string]string{"A": "1", "B": "2", "C": "3"}

	kvs := r.Apply(secrets)

	if len(kvs) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(kvs))
	}
	expected := []string{"C", "A", "B"}
	for i, kv := range kvs {
		if kv.Key != expected[i] {
			t.Errorf("position %d: want %s, got %s", i, expected[i], kv.Key)
		}
	}
}

func TestApply_RemainingKeysSortedWithFallback(t *testing.T) {
	r, _ := envreorder.New([]string{"Z"}, true)
	secrets := map[string]string{"Z": "last", "A": "first", "M": "middle"}

	kvs := r.Apply(secrets)

	if kvs[0].Key != "Z" {
		t.Errorf("first key should be Z, got %s", kvs[0].Key)
	}
	if kvs[1].Key != "A" {
		t.Errorf("second key should be A (sorted), got %s", kvs[1].Key)
	}
	if kvs[2].Key != "M" {
		t.Errorf("third key should be M (sorted), got %s", kvs[2].Key)
	}
}

func TestApply_OrderedKeyMissing_Skipped(t *testing.T) {
	r, _ := envreorder.New([]string{"X", "Y"}, true)
	secrets := map[string]string{"Y": "yes", "A": "alpha"}

	kvs := r.Apply(secrets)

	if kvs[0].Key != "Y" {
		t.Errorf("expected Y first, got %s", kvs[0].Key)
	}
	if kvs[1].Key != "A" {
		t.Errorf("expected A second, got %s", kvs[1].Key)
	}
}

func TestKeys_ReturnsOnlyKeys(t *testing.T) {
	r, _ := envreorder.New([]string{"B", "A"}, false)
	secrets := map[string]string{"A": "1", "B": "2"}

	keys := r.Keys(secrets)

	if len(keys) != 2 || keys[0] != "B" || keys[1] != "A" {
		t.Errorf("unexpected keys: %v", keys)
	}
}

func TestApply_EmptySecrets(t *testing.T) {
	r, _ := envreorder.New([]string{"A"}, true)
	kvs := r.Apply(map[string]string{})
	if len(kvs) != 0 {
		t.Errorf("expected empty result, got %d entries", len(kvs))
	}
}
