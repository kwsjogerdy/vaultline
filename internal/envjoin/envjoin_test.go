package envjoin_test

import (
	"strings"
	"testing"

	"github.com/vaultline/vaultline/internal/envjoin"
)

func TestNew_UnknownMode_ReturnsError(t *testing.T) {
	_, err := envjoin.New("invalid", ",")
	if err == nil {
		t.Fatal("expected error for unknown mode")
	}
}

func TestJoin_LeftOnly(t *testing.T) {
	j, _ := envjoin.New(envjoin.ModeConcat, ",")
	left := map[string]string{"A": "1"}
	right := map[string]string{}
	out, err := j.Join(left, right)
	if err != nil {
		t.Fatal(err)
	}
	if out["A"] != "1" {
		t.Errorf("expected A=1, got %q", out["A"])
	}
}

func TestJoin_RightOnly(t *testing.T) {
	j, _ := envjoin.New(envjoin.ModeRight, "")
	left := map[string]string{}
	right := map[string]string{"B": "2"}
	out, err := j.Join(left, right)
	if err != nil {
		t.Fatal(err)
	}
	if out["B"] != "2" {
		t.Errorf("expected B=2, got %q", out["B"])
	}
}

func TestJoin_ConcatMode(t *testing.T) {
	j, _ := envjoin.New(envjoin.ModeConcat, ":")
	left := map[string]string{"URL": "http://host"}
	right := map[string]string{"URL": "8080"}
	out, err := j.Join(left, right)
	if err != nil {
		t.Fatal(err)
	}
	if out["URL"] != "http://host:8080" {
		t.Errorf("unexpected URL value: %q", out["URL"])
	}
}

func TestJoin_LeftMode_KeepsLeft(t *testing.T) {
	j, _ := envjoin.New(envjoin.ModeLeft, "")
	left := map[string]string{"KEY": "original"}
	right := map[string]string{"KEY": "override"}
	out, err := j.Join(left, right)
	if err != nil {
		t.Fatal(err)
	}
	if out["KEY"] != "original" {
		t.Errorf("expected original, got %q", out["KEY"])
	}
}

func TestJoin_RightMode_KeepsRight(t *testing.T) {
	j, _ := envjoin.New(envjoin.ModeRight, "")
	left := map[string]string{"KEY": "original"}
	right := map[string]string{"KEY": "override"}
	out, err := j.Join(left, right)
	if err != nil {
		t.Fatal(err)
	}
	if out["KEY"] != "override" {
		t.Errorf("expected override, got %q", out["KEY"])
	}
}

func TestJoin_NilLeft_ReturnsError(t *testing.T) {
	j, _ := envjoin.New(envjoin.ModeLeft, "")
	_, err := j.Join(nil, map[string]string{})
	if err == nil {
		t.Fatal("expected error for nil left map")
	}
}

func TestJoin_NilRight_ReturnsError(t *testing.T) {
	j, _ := envjoin.New(envjoin.ModeLeft, "")
	_, err := j.Join(map[string]string{}, nil)
	if err == nil {
		t.Fatal("expected error for nil right map")
	}
}

func TestSummary_ContainsExpectedCounts(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"C": "3"}
	result := map[string]string{"A": "1", "B": "2", "C": "3"}
	s := envjoin.Summary(left, right, result)
	if !strings.Contains(s, "left=2") || !strings.Contains(s, "right=1") || !strings.Contains(s, "result=3") {
		t.Errorf("unexpected summary: %q", s)
	}
}
