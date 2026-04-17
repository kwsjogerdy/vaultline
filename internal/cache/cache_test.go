package cache_test

import (
	"os"
	"testing"
	"time"

	"github.com/vaultline/vaultline/internal/cache"
)

func TestGet_MissWhenEmpty(t *testing.T) {
	dir := t.TempDir()
	c := cache.New(dir, 5*time.Minute)
	_, ok := c.Get("secret/myapp")
	if ok {
		t.Fatal("expected cache miss, got hit")
	}
}

func TestSetAndGet_HitWithinTTL(t *testing.T) {
	dir := t.TempDir()
	c := cache.New(dir, 5*time.Minute)
	secrets := map[string]string{"DB_PASS": "hunter2", "API_KEY": "abc123"}

	if err := c.Set("secret/myapp", secrets); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	got, ok := c.Get("secret/myapp")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if got["DB_PASS"] != "hunter2" {
		t.Errorf("unexpected value: %s", got["DB_PASS"])
	}
}

func TestGet_MissAfterTTL(t *testing.T) {
	dir := t.TempDir()
	c := cache.New(dir, 1*time.Millisecond)
	secrets := map[string]string{"KEY": "val"}

	if err := c.Set("secret/myapp", secrets); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	time.Sleep(5 * time.Millisecond)

	_, ok := c.Get("secret/myapp")
	if ok {
		t.Fatal("expected cache miss after TTL expiry")
	}
}

func TestInvalidate_RemovesEntry(t *testing.T) {
	dir := t.TempDir()
	c := cache.New(dir, 5*time.Minute)
	secrets := map[string]string{"X": "y"}

	_ = c.Set("secret/myapp", secrets)
	if err := c.Invalidate("secret/myapp"); err != nil {
		t.Fatalf("Invalidate failed: %v", err)
	}
	_, ok := c.Get("secret/myapp")
	if ok {
		t.Fatal("expected miss after invalidation")
	}
}

func TestInvalidate_NoopWhenMissing(t *testing.T) {
	dir := t.TempDir()
	c := cache.New(dir, 5*time.Minute)
	if err := c.Invalidate("secret/nonexistent"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestSet_CreatesDir(t *testing.T) {
	dir := t.TempDir()
	subdir := dir + "/nested/cache"
	c := cache.New(subdir, 5*time.Minute)

	if err := c.Set("secret/app", map[string]string{"A": "B"}); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	if _, err := os.Stat(subdir); err != nil {
		t.Fatalf("expected dir to be created: %v", err)
	}
}
