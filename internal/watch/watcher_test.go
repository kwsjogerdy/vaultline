package watch_test

import (
	"context"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourorg/vaultline/internal/watch"
)

func writeTmp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "vaultline-watch-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(f.Name(), []byte(content), 0600)
	return f.Name()
}

func TestWatcher_DetectsChange(t *testing.T) {
	path := writeTmp(t, "initial")

	var calls atomic.Int32
	w := watch.New(path, 20*time.Millisecond, func() error {
		calls.Add(1)
		return nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	go w.Run(ctx) //nolint:errcheck

	time.Sleep(50 * time.Millisecond)
	_ = os.WriteFile(path, []byte("changed"), 0600)
	time.Sleep(100 * time.Millisecond)

	if calls.Load() == 0 {
		t.Error("expected sync fn to be called after file change")
	}
}

func TestWatcher_NoCallWithoutChange(t *testing.T) {
	path := writeTmp(t, "stable")

	var calls atomic.Int32
	w := watch.New(path, 20*time.Millisecond, func() error {
		calls.Add(1)
		return nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	go w.Run(ctx) //nolint:errcheck
	time.Sleep(160 * time.Millisecond)

	if calls.Load() != 0 {
		t.Errorf("expected no calls, got %d", calls.Load())
	}
}

func TestWatcher_MissingFile(t *testing.T) {
	w := watch.New("/nonexistent/file.yaml", 20*time.Millisecond, func() error { return nil })
	ctx := context.Background()
	if err := w.Run(ctx); err == nil {
		t.Error("expected error for missing file")
	}
}
