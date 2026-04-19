package envpromote

import (
	"testing"
)

func TestPromote_CopiesAllKeys(t *testing.T) {
	p := New("dev", "staging")
	src := map[string]string{"DB_HOST": "localhost", "API_KEY": "abc"}
	dst := map[string]string{}
	out, res := p.Promote(src, dst, nil)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if len(res.Copied) != 2 {
		t.Errorf("expected 2 copied, got %d", len(res.Copied))
	}
}

func TestPromote_OnlyKeys(t *testing.T) {
	p := New("dev", "staging")
	src := map[string]string{"DB_HOST": "localhost", "API_KEY": "abc", "SECRET": "x"}
	out, res := p.Promote(src, map[string]string{}, []string{"DB_HOST"})
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in output")
	}
	if _, ok := out["API_KEY"]; ok {
		t.Error("API_KEY should have been skipped")
	}
	if len(res.Skipped) != 2 {
		t.Errorf("expected 2 skipped, got %d", len(res.Skipped))
	}
}

func TestPromote_StripPrefix(t *testing.T) {
	p := New("dev", "staging").WithPrefix("DEV_")
	src := map[string]string{"DEV_DB_HOST": "localhost", "GLOBAL": "yes"}
	out, _ := p.Promote(src, map[string]string{}, nil)
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected DEV_ prefix stripped")
	}
	if _, ok := out["GLOBAL"]; !ok {
		t.Error("expected GLOBAL key preserved")
	}
}

func TestPromote_MergesIntoDst(t *testing.T) {
	p := New("dev", "staging")
	src := map[string]string{"NEW_KEY": "new"}
	dst := map[string]string{"EXISTING": "keep"}
	out, _ := p.Promote(src, dst, nil)
	if out["EXISTING"] != "keep" {
		t.Error("existing dst key should be preserved")
	}
	if out["NEW_KEY"] != "new" {
		t.Error("new key should be added")
	}
}

func TestPromote_EmptySrc(t *testing.T) {
	p := New("dev", "staging")
	out, res := p.Promote(map[string]string{}, map[string]string{"A": "1"}, nil)
	if out["A"] != "1" {
		t.Error("dst should be unchanged")
	}
	if len(res.Copied) != 0 {
		t.Error("nothing should be copied")
	}
}
