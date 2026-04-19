package envnormalize

import (
	"testing"
)

func TestApply_UppercaseKeys(t *testing.T) {
	n := New(Options{UppercaseKeys: true})
	out := n.Apply(map[string]string{"db_host": "localhost"})
	if out["DB_HOST"] != "localhost" {
		t.Fatalf("expected DB_HOST, got %v", out)
	}
}

func TestApply_TrimValues(t *testing.T) {
	n := New(Options{TrimValues: true})
	out := n.Apply(map[string]string{"KEY": "  value  "})
	if out["KEY"] != "value" {
		t.Fatalf("expected trimmed value, got %q", out["KEY"])
	}
}

func TestApply_ReplaceHyphens(t *testing.T) {
	n := New(Options{ReplaceHyphens: true})
	out := n.Apply(map[string]string{"my-key": "v"})
	if _, ok := out["my_key"]; !ok {
		t.Fatalf("expected my_key in output, got %v", out)
	}
}

func TestApply_StripNonAlnum(t *testing.T) {
	n := New(Options{StripNonAlnum: true})
	out := n.Apply(map[string]string{"KEY@NAME!": "v"})
	if _, ok := out["KEYNAME"]; !ok {
		t.Fatalf("expected KEYNAME, got %v", out)
	}
}

func TestApply_EmptyKeyDropped(t *testing.T) {
	n := New(Options{StripNonAlnum: true})
	out := n.Apply(map[string]string{"!!!": "v"})
	if len(out) != 0 {
		t.Fatalf("expected empty map, got %v", out)
	}
}

func TestApply_DefaultOptions(t *testing.T) {
	n := New(DefaultOptions())
	out := n.Apply(map[string]string{"my-secret": "  hello  "})
	if out["MY_SECRET"] != "hello" {
		t.Fatalf("unexpected result: %v", out)
	}
}

func TestApply_EmptySecrets(t *testing.T) {
	n := New(DefaultOptions())
	out := n.Apply(map[string]string{})
	if len(out) != 0 {
		t.Fatalf("expected empty output")
	}
}
