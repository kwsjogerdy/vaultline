package envdrift

import (
	"testing"
)

func TestDetect_AllMatch(t *testing.T) {
	local := map[string]string{"A": "1", "B": "2"}
	vault := map[string]string{"A": "1", "B": "2"}
	r := Detect(local, vault)
	if r.HasDrift() {
		t.Fatal("expected no drift")
	}
}

func TestDetect_Drifted(t *testing.T) {
	local := map[string]string{"A": "old"}
	vault := map[string]string{"A": "new"}
	r := Detect(local, vault)
	if !r.HasDrift() {
		t.Fatal("expected drift")
	}
	if r.Entries[0].Status != StatusDrifted {
		t.Fatalf("expected drifted, got %s", r.Entries[0].Status)
	}
}

func TestDetect_Missing(t *testing.T) {
	local := map[string]string{}
	vault := map[string]string{"SECRET": "val"}
	r := Detect(local, vault)
	if r.Entries[0].Status != StatusMissing {
		t.Fatalf("expected missing, got %s", r.Entries[0].Status)
	}
}

func TestDetect_Extra(t *testing.T) {
	local := map[string]string{"EXTRA": "val"}
	vault := map[string]string{}
	r := Detect(local, vault)
	if r.Entries[0].Status != StatusExtra {
		t.Fatalf("expected extra, got %s", r.Entries[0].Status)
	}
}

func TestReport_Summary(t *testing.T) {
	local := map[string]string{"A": "1", "B": "old", "C": "extra"}
	vault := map[string]string{"A": "1", "B": "new", "D": "missing"}
	r := Detect(local, vault)
	s := r.Summary()
	if s == "" {
		t.Fatal("expected non-empty summary")
	}
}

func TestStatus_String(t *testing.T) {
	cases := []struct {
		s    Status
		want string
	}{
		{StatusMatch, "match"},
		{StatusDrifted, "drifted"},
		{StatusMissing, "missing"},
		{StatusExtra, "extra"},
		{Status(99), "unknown"},
	}
	for _, c := range cases {
		if got := c.s.String(); got != c.want {
			t.Errorf("Status(%d).String() = %q, want %q", c.s, got, c.want)
		}
	}
}
