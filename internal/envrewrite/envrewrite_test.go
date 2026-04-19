package envrewrite_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envrewrite"
)

func TestApply_RewriteKeys(t *testing.T) {
	r := envrewrite.New([]envrewrite.Rule{
		{Find: "DEV_", Replace: "PROD_", Target: "key"},
	})
	out, err := r.Apply(map[string]string{"DEV_DB": "localhost", "DEV_PORT": "5432"})
	if err != nil {
		t.Fatal(err)
	}
	if out["PROD_DB"] != "localhost" {
		t.Errorf("expected PROD_DB=localhost, got %v", out)
	}
	if out["PROD_PORT"] != "5432" {
		t.Errorf("expected PROD_PORT=5432, got %v", out)
	}
}

func TestApply_RewriteValues(t *testing.T) {
	r := envrewrite.New([]envrewrite.Rule{
		{Find: "localhost", Replace: "db.prod", Target: "value"},
	})
	out, err := r.Apply(map[string]string{"DB_HOST": "localhost", "APP": "myapp"})
	if err != nil {
		t.Fatal(err)
	}
	if out["DB_HOST"] != "db.prod" {
		t.Errorf("expected db.prod, got %s", out["DB_HOST"])
	}
	if out["APP"] != "myapp" {
		t.Errorf("APP should be unchanged")
	}
}

func TestApply_RewriteBoth(t *testing.T) {
	r := envrewrite.New([]envrewrite.Rule{
		{Find: "OLD", Replace: "NEW", Target: "both"},
	})
	out, err := r.Apply(map[string]string{"OLD_KEY": "OLD_VALUE"})
	if err != nil {
		t.Fatal(err)
	}
	if out["NEW_KEY"] != "NEW_VALUE" {
		t.Errorf("expected NEW_KEY=NEW_VALUE, got %v", out)
	}
}

func TestApply_UnknownTarget(t *testing.T) {
	r := envrewrite.New([]envrewrite.Rule{
		{Find: "X", Replace: "Y", Target: "invalid"},
	})
	_, err := r.Apply(map[string]string{"X": "1"})
	if err == nil {
		t.Fatal("expected error for unknown target")
	}
}

func TestApply_EmptyFind_ReturnsError(t *testing.T) {
	r := envrewrite.New([]envrewrite.Rule{
		{Find: "", Replace: "Y", Target: "key"},
	})
	_, err := r.Apply(map[string]string{"A": "1"})
	if err == nil {
		t.Fatal("expected error for empty find")
	}
}

func TestApply_NoRules_ReturnsUnchanged(t *testing.T) {
	r := envrewrite.New(nil)
	in := map[string]string{"FOO": "bar"}
	out, err := r.Apply(in)
	if err != nil {
		t.Fatal(err)
	}
	if out["FOO"] != "bar" {
		t.Errorf("expected unchanged map")
	}
}
