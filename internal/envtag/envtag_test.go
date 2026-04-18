package envtag_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envtag"
)

func TestAdd_GroupsByLabel(t *testing.T) {
	tr := envtag.New()
	tr.Add("prod", "DB_HOST")
	tr.Add("prod", "DB_PASS")
	tr.Add("dev", "DEBUG")

	tag, err := tr.Get("prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tag.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(tag.Keys))
	}
}

func TestGet_NotFound(t *testing.T) {
	tr := envtag.New()
	_, err := tr.Get("missing")
	if err == nil {
		t.Error("expected error for missing label")
	}
}

func TestInferAndAdd_KnownPrefix(t *testing.T) {
	tr := envtag.New()
	tr.InferAndAdd("PROD_API_KEY")
	tr.InferAndAdd("DEV_SECRET")
	tr.InferAndAdd("UNTAGGED_VAR")

	prod, err := tr.Get("prod")
	if err != nil || len(prod.Keys) != 1 {
		t.Errorf("expected 1 prod key, got err=%v", err)
	}
	dev, err := tr.Get("dev")
	if err != nil || len(dev.Keys) != 1 {
		t.Errorf("expected 1 dev key, got err=%v", err)
	}
	def, err := tr.Get("default")
	if err != nil || len(def.Keys) != 1 {
		t.Errorf("expected 1 default key, got err=%v", err)
	}
}

func TestFilter_ReturnsMatchingSecrets(t *testing.T) {
	tr := envtag.New()
	tr.Add("prod", "API_KEY")
	tr.Add("prod", "DB_PASS")

	secrets := map[string]string{
		"API_KEY": "abc123",
		"DB_PASS": "secret",
		"OTHER":   "ignore",
	}

	result, err := tr.Filter("prod", secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 filtered secrets, got %d", len(result))
	}
	if _, ok := result["OTHER"]; ok {
		t.Error("OTHER should not be in filtered result")
	}
}

func TestFilter_LabelNotFound(t *testing.T) {
	tr := envtag.New()
	_, err := tr.Filter("staging", map[string]string{})
	if err == nil {
		t.Error("expected error for unknown label")
	}
}

func TestAll_ReturnsAllTags(t *testing.T) {
	tr := envtag.New()
	tr.Add("prod", "X")
	tr.Add("dev", "Y")
	all := tr.All()
	if len(all) != 2 {
		t.Errorf("expected 2 tags, got %d", len(all))
	}
}
