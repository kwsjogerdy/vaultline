package prompt

import (
	"bytes"
	"strings"
	"testing"
)

func TestConfirm_Yes(t *testing.T) {
	for _, input := range []string{"y\n", "Y\n", "yes\n", "YES\n"} {
		in := strings.NewReader(input)
		out := &bytes.Buffer{}
		p := NewWithIO(in, out)
		ok, err := p.Confirm("Continue?")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ok {
			t.Errorf("expected true for input %q", input)
		}
	}
}

func TestConfirm_No(t *testing.T) {
	for _, input := range []string{"n\n", "N\n", "no\n", "\n", "maybe\n"} {
		in := strings.NewReader(input)
		out := &bytes.Buffer{}
		p := NewWithIO(in, out)
		ok, err := p.Confirm("Continue?")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ok {
			t.Errorf("expected false for input %q", input)
		}
	}
}

func TestConfirm_EOF(t *testing.T) {
	in := strings.NewReader("")
	out := &bytes.Buffer{}
	p := NewWithIO(in, out)
	ok, err := p.Confirm("Continue?")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected false on EOF")
	}
}

func TestConfirmDiff_OutputFormat(t *testing.T) {
	in := strings.NewReader("y\n")
	out := &bytes.Buffer{}
	p := NewWithIO(in, out)
	ok, err := p.ConfirmDiff(3, 1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected true")
	}
	output := out.String()
	for _, want := range []string{"+ 3 added", "- 1 removed", "~ 2 changed"} {
		if !strings.Contains(output, want) {
			t.Errorf("output missing %q", want)
		}
	}
}
