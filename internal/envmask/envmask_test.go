package envmask_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envmask"
)

func TestApply_NoRules_ReturnsUnchanged(t *testing.T) {
	m, err := envmask.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	secrets := map[string]string{"DB_HOST": "localhost", "PORT": "5432"}
	out := m.Apply(secrets)
	if out["DB_HOST"] != "localhost" || out["PORT"] != "5432" {
		t.Errorf("expected unchanged values, got %v", out)
	}
}

func TestApply_FullRedact(t *testing.T) {
	m, _ := envmask.New([]envmask.Rule{{Pattern: "password", Reveal: 0}})
	out := m.Apply(map[string]string{"DB_PASSWORD": "supersecret"})
	if out["DB_PASSWORD"] != "***" {
		t.Errorf("expected ***, got %s", out["DB_PASSWORD"])
	}
}

func TestApply_RevealSuffix(t *testing.T) {
	m, _ := envmask.New([]envmask.Rule{{Pattern: "token", Reveal: 4}})
	out := m.Apply(map[string]string{"API_TOKEN": "abcdef1234"})
	if out["API_TOKEN"] != "***1234" {
		t.Errorf("expected ***1234, got %s", out["API_TOKEN"])
	}
}

func TestApply_RevealExceedsLength(t *testing.T) {
	m, _ := envmask.New([]envmask.Rule{{Pattern: "secret", Reveal: 100}})
	out := m.Apply(map[string]string{"MY_SECRET": "short"})
	if out["MY_SECRET"] != "***" {
		t.Errorf("expected ***, got %s", out["MY_SECRET"])
	}
}

func TestApply_CaseInsensitivePattern(t *testing.T) {
	m, _ := envmask.New([]envmask.Rule{{Pattern: "KEY", Reveal: 0}})
	out := m.Apply(map[string]string{"aws_key": "value123"})
	if out["aws_key"] != "***" {
		t.Errorf("expected ***, got %s", out["aws_key"])
	}
}

func TestApply_MultipleRules_FirstMatchWins(t *testing.T) {
	m, _ := envmask.New([]envmask.Rule{
		{Pattern: "secret", Reveal: 0},
		{Pattern: "secret", Reveal: 4},
	})
	out := m.Apply(map[string]string{"MY_SECRET": "abcdefgh"})
	if out["MY_SECRET"] != "***" {
		t.Errorf("expected first rule to win, got %s", out["MY_SECRET"])
	}
}

func TestNew_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := envmask.New([]envmask.Rule{{Pattern: "[invalid"}})
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestMaskValue_DirectCall(t *testing.T) {
	m, _ := envmask.New([]envmask.Rule{{Pattern: "pass", Reveal: 2}})
	result := m.MaskValue("DB_PASSWORD", "mypassword")
	if result != "***rd" {
		t.Errorf("expected ***rd, got %s", result)
	}
}
