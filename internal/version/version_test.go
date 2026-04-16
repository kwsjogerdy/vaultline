package version

import (
	"runtime"
	"strings"
	"testing"
)

func TestGet_Defaults(t *testing.T) {
	info := Get()
	if info.Version != "dev" {
		t.Errorf("expected default version 'dev', got %q", info.Version)
	}
	if info.Commit != "none" {
		t.Errorf("expected default commit 'none', got %q", info.Commit)
	}
	if info.BuildDate != "unknown" {
		t.Errorf("expected default build date 'unknown', got %q", info.BuildDate)
	}
	if info.GoVersion != runtime.Version() {
		t.Errorf("expected go version %q, got %q", runtime.Version(), info.GoVersion)
	}
}

func TestInfo_String_ContainsFields(t *testing.T) {
	info := Info{
		Version:   "1.2.3",
		Commit:    "abc1234",
		BuildDate: "2024-01-15",
		GoVersion: "go1.22.0",
	}
	s := info.String()
	for _, want := range []string{"1.2.3", "abc1234", "2024-01-15", "go1.22.0", "vaultline"} {
		if !strings.Contains(s, want) {
			t.Errorf("expected string to contain %q, got: %s", want, s)
		}
	}
}

func TestGet_OverrideVars(t *testing.T) {
	origVersion := Version
	origCommit := Commit
	defer func() {
		Version = origVersion
		Commit = origCommit
	}()

	Version = "0.9.0"
	Commit = "deadbeef"

	info := Get()
	if info.Version != "0.9.0" {
		t.Errorf("expected '0.9.0', got %q", info.Version)
	}
	if info.Commit != "deadbeef" {
		t.Errorf("expected 'deadbeef', got %q", info.Commit)
	}
}
