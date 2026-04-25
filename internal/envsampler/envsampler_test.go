package envsampler_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envsampler"
)

var testSecrets = map[string]string{
	"APP_HOST":    "localhost",
	"APP_PORT":    "8080",
	"DB_HOST":     "db.local",
	"DB_PASSWORD": "secret",
	"DB_USER":     "admin",
	"REDIS_URL":   "redis://localhost",
}

func TestNew_InvalidN(t *testing.T) {
	_, err := envsampler.New(envsampler.Options{N: 0})
	if err == nil {
		t.Fatal("expected error for N=0")
	}
}

func TestSample_ReturnsAtMostN(t *testing.T) {
	s, err := envsampler.New(envsampler.Options{N: 2, Seed: 42})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := s.Sample(testSecrets)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}

func TestSample_NGreaterThanTotal(t *testing.T) {
	s, _ := envsampler.New(envsampler.Options{N: 100, Seed: 1})
	out := s.Sample(testSecrets)
	if len(out) != len(testSecrets) {
		t.Fatalf("expected %d keys, got %d", len(testSecrets), len(out))
	}
}

func TestSample_DeterministicWithSeed(t *testing.T) {
	s1, _ := envsampler.New(envsampler.Options{N: 3, Seed: 99})
	s2, _ := envsampler.New(envsampler.Options{N: 3, Seed: 99})
	out1 := s1.Keys(testSecrets)
	out2 := s2.Keys(testSecrets)
	if len(out1) != len(out2) {
		t.Fatal("lengths differ")
	}
	for i := range out1 {
		if out1[i] != out2[i] {
			t.Fatalf("key mismatch at %d: %s vs %s", i, out1[i], out2[i])
		}
	}
}

func TestSample_PrefixFilter(t *testing.T) {
	s, _ := envsampler.New(envsampler.Options{N: 10, Seed: 7, Prefix: "DB_"})
	out := s.Sample(testSecrets)
	for k := range out {
		if len(k) < 3 || k[:3] != "DB_" {
			t.Fatalf("unexpected key without DB_ prefix: %s", k)
		}
	}
	if len(out) != 3 {
		t.Fatalf("expected 3 DB_ keys, got %d", len(out))
	}
}

func TestSample_EmptySecrets(t *testing.T) {
	s, _ := envsampler.New(envsampler.Options{N: 5, Seed: 1})
	out := s.Sample(map[string]string{})
	if len(out) != 0 {
		t.Fatal("expected empty result")
	}
}

func TestKeys_ReturnsSorted(t *testing.T) {
	s, _ := envsampler.New(envsampler.Options{N: 10, Seed: 3})
	keys := s.Keys(testSecrets)
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Fatalf("keys not sorted at index %d: %s < %s", i, keys[i], keys[i-1])
		}
	}
}
