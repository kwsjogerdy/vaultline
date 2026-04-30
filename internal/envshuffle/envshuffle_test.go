package envshuffle_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envshuffle"
)

var sample = map[string]string{
	"APP_HOST": "localhost",
	"APP_PORT": "8080",
	"DB_HOST":  "db.local",
	"DB_PASS":  "secret",
	"LOG_LEVEL": "info",
}

func TestApply_ReturnsSameCount(t *testing.T) {
	s := envshuffle.New()
	entries, err := s.Apply(sample)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != len(sample) {
		t.Fatalf("expected %d entries, got %d", len(sample), len(entries))
	}
}

func TestApply_DeterministicWithSeed(t *testing.T) {
	s1 := envshuffle.New(envshuffle.WithSeed(42))
	s2 := envshuffle.New(envshuffle.WithSeed(42))

	e1, _ := s1.Apply(sample)
	e2, _ := s2.Apply(sample)

	for i := range e1 {
		if e1[i].Key != e2[i].Key {
			t.Fatalf("index %d: expected key %q, got %q", i, e1[i].Key, e2[i].Key)
		}
	}
}

func TestApply_DifferentSeedsDiffer(t *testing.T) {
	s1 := envshuffle.New(envshuffle.WithSeed(1))
	s2 := envshuffle.New(envshuffle.WithSeed(99999))

	e1, _ := s1.Apply(sample)
	e2, _ := s2.Apply(sample)

	same := true
	for i := range e1 {
		if e1[i].Key != e2[i].Key {
			same = false
			break
		}
	}
	if same {
		t.Error("expected different order for different seeds")
	}
}

func TestApply_PrefixFilter(t *testing.T) {
	s := envshuffle.New(envshuffle.WithPrefix("APP_"))
	entries, err := s.Apply(sample)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	for _, e := range entries {
		if len(e.Key) < 4 || e.Key[:4] != "APP_" {
			t.Errorf("unexpected key %q does not match prefix APP_", e.Key)
		}
	}
}

func TestApply_Reverse_SortedDescending(t *testing.T) {
	s := envshuffle.New(envshuffle.WithReverse())
	entries, err := s.Apply(sample)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 1; i < len(entries); i++ {
		if entries[i-1].Key < entries[i].Key {
			t.Errorf("expected descending order, got %q before %q", entries[i-1].Key, entries[i].Key)
		}
	}
}

func TestApply_NilMap_ReturnsError(t *testing.T) {
	s := envshuffle.New()
	_, err := s.Apply(nil)
	if err == nil {
		t.Fatal("expected error for nil map")
	}
}
