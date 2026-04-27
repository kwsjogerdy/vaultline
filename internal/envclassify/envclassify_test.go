package envclassify_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envclassify"
)

func TestClassify_Credential(t *testing.T) {
	c := envclassify.New()
	secrets := map[string]string{"DB_PASSWORD": "s3cr3t", "API_TOKEN": "tok"}
	res := c.Classify(secrets)
	for _, r := range res {
		if r.Category != envclassify.CategoryCredential {
			t.Errorf("key %q: want credential, got %q", r.Key, r.Category)
		}
	}
}

func TestClassify_Endpoint(t *testing.T) {
	c := envclassify.New()
	secrets := map[string]string{"DATABASE_URL": "postgres://localhost", "REDIS_HOST": "127.0.0.1"}
	res := c.Classify(secrets)
	for _, r := range res {
		if r.Category != envclassify.CategoryEndpoint {
			t.Errorf("key %q: want endpoint, got %q", r.Key, r.Category)
		}
	}
}

func TestClassify_FeatureFlag(t *testing.T) {
	c := envclassify.New()
	secrets := map[string]string{"ENABLE_DARK_MODE": "true", "FEATURE_BETA": "1"}
	res := c.Classify(secrets)
	for _, r := range res {
		if r.Category != envclassify.CategoryFeatureFlag {
			t.Errorf("key %q: want feature_flag, got %q", r.Key, r.Category)
		}
	}
}

func TestClassify_Identifier(t *testing.T) {
	c := envclassify.New()
	secrets := map[string]string{"PROJECT_NAME": "vaultline", "REGION": "us-east-1"}
	res := c.Classify(secrets)
	for _, r := range res {
		if r.Category != envclassify.CategoryIdentifier {
			t.Errorf("key %q: want identifier, got %q", r.Key, r.Category)
		}
	}
}

func TestClassify_Unknown(t *testing.T) {
	c := envclassify.New()
	secrets := map[string]string{"FOOBAR": "baz"}
	res := c.Classify(secrets)
	if len(res) != 1 || res[0].Category != envclassify.CategoryUnknown {
		t.Fatalf("expected unknown, got %q", res[0].Category)
	}
}

func TestClassify_SortedOutput(t *testing.T) {
	c := envclassify.New()
	secrets := map[string]string{"Z_TOKEN": "a", "A_URL": "b", "M_ID": "c"}
	res := c.Classify(secrets)
	if res[0].Key != "A_URL" || res[1].Key != "M_ID" || res[2].Key != "Z_TOKEN" {
		t.Errorf("results not sorted: %v", res)
	}
}

func TestByCategory_Groups(t *testing.T) {
	c := envclassify.New()
	secrets := map[string]string{"DB_PASSWORD": "x", "API_URL": "y", "FOOBAR": "z"}
	res := c.Classify(secrets)
	grouped := envclassify.ByCategory(res)
	if len(grouped[envclassify.CategoryCredential]) != 1 {
		t.Error("expected 1 credential")
	}
	if len(grouped[envclassify.CategoryEndpoint]) != 1 {
		t.Error("expected 1 endpoint")
	}
	if len(grouped[envclassify.CategoryUnknown]) != 1 {
		t.Error("expected 1 unknown")
	}
}
