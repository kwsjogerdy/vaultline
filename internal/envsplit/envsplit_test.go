package envsplit_test

import (
	"testing"

	"github.com/vaultline/vaultline/internal/envsplit"
)

func TestNew_EmptyBucketName_ReturnsError(t *testing.T) {
	_, err := envsplit.New([]envsplit.Rule{{Bucket: "", Prefix: "APP_"}})
	if err == nil {
		t.Fatal("expected error for empty bucket name")
	}
}

func TestApply_RoutesToCorrectBuckets(t *testing.T) {
	s, err := envsplit.New([]envsplit.Rule{
		{Bucket: "app", Prefix: "APP_"},
		{Bucket: "db", Prefix: "DB_"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secrets := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_URL":   "postgres://",
		"OTHER":    "value",
	}

	res := s.Apply(secrets)

	if res.Buckets["app"]["APP_HOST"] != "localhost" {
		t.Error("expected APP_HOST in app bucket")
	}
	if res.Buckets["app"]["APP_PORT"] != "8080" {
		t.Error("expected APP_PORT in app bucket")
	}
	if res.Buckets["db"]["DB_URL"] != "postgres://" {
		t.Error("expected DB_URL in db bucket")
	}
	if res.Unrouted["OTHER"] != "value" {
		t.Error("expected OTHER in unrouted")
	}
}

func TestApply_EmptySecrets(t *testing.T) {
	s, _ := envsplit.New([]envsplit.Rule{{Bucket: "app", Prefix: "APP_"}})
	res := s.Apply(map[string]string{})
	if len(res.Buckets) != 0 || len(res.Unrouted) != 0 {
		t.Error("expected empty result for empty secrets")
	}
}

func TestApply_FirstRuleWins(t *testing.T) {
	s, _ := envsplit.New([]envsplit.Rule{
		{Bucket: "first", Prefix: "APP_"},
		{Bucket: "second", Prefix: "APP_"},
	})
	res := s.Apply(map[string]string{"APP_KEY": "v"})
	if _, ok := res.Buckets["first"]["APP_KEY"]; !ok {
		t.Error("expected first rule to win")
	}
	if _, ok := res.Buckets["second"]; ok {
		t.Error("second bucket should be empty")
	}
}

func TestApply_EmptyPrefixCatchAll(t *testing.T) {
	s, _ := envsplit.New([]envsplit.Rule{
		{Bucket: "specific", Prefix: "DB_"},
		{Bucket: "all", Prefix: ""},
	})
	res := s.Apply(map[string]string{
		"DB_HOST": "localhost",
		"OTHER":   "value",
	})
	if res.Buckets["specific"]["DB_HOST"] != "localhost" {
		t.Error("expected DB_HOST in specific bucket")
	}
	if res.Buckets["all"]["OTHER"] != "value" {
		t.Error("expected OTHER in all bucket")
	}
	if len(res.Unrouted) != 0 {
		t.Error("expected no unrouted keys when catch-all rule present")
	}
}

func TestBucketNames_SortedOrder(t *testing.T) {
	s, _ := envsplit.New([]envsplit.Rule{
		{Bucket: "zebra", Prefix: "Z_"},
		{Bucket: "alpha", Prefix: "A_"},
	})
	res := s.Apply(map[string]string{"Z_KEY": "1", "A_KEY": "2"})
	names := envsplit.BucketNames(res)
	if len(names) != 2 || names[0] != "alpha" || names[1] != "zebra" {
		t.Errorf("expected sorted bucket names, got %v", names)
	}
}
