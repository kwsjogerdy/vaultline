package envmerge_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envmerge"
)

func TestMerge_NoConflict(t *testing.T) {
	m := envmerge.New(envmerge.StrategyFirst)
	a := map[string]string{"A": "1", "B": "2"}
	b := map[string]string{"C": "3"}
	res, err := m.Merge(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Secrets) != 3 {
		t.Errorf("expected 3 keys, got %d", len(res.Secrets))
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts")
	}
}

func TestMerge_StrategyFirst_KeepsFirst(t *testing.T) {
	m := envmerge.New(envmerge.StrategyFirst)
	a := map[string]string{"KEY": "original"}
	b := map[string]string{"KEY": "override"}
	res, err := m.Merge(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Secrets["KEY"] != "original" {
		t.Errorf("expected 'original', got %q", res.Secrets["KEY"])
	}
	if len(res.Conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d", len(res.Conflicts))
	}
}

func TestMerge_StrategyLast_KeepsLast(t *testing.T) {
	m := envmerge.New(envmerge.StrategyLast)
	a := map[string]string{"KEY": "original"}
	b := map[string]string{"KEY": "override"}
	res, err := m.Merge(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Secrets["KEY"] != "override" {
		t.Errorf("expected 'override', got %q", res.Secrets["KEY"])
	}
}

func TestMerge_StrategyError_ReturnsError(t *testing.T) {
	m := envmerge.New(envmerge.StrategyError)
	a := map[string]string{"KEY": "v1"}
	b := map[string]string{"KEY": "v2"}
	_, err := m.Merge(a, b)
	if err == nil {
		t.Fatal("expected error on conflict, got nil")
	}
}

func TestMerge_EmptySources(t *testing.T) {
	m := envmerge.New(envmerge.StrategyFirst)
	res, err := m.Merge()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Secrets) != 0 {
		t.Errorf("expected empty result")
	}
}
