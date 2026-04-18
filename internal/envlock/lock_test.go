package envlock_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vaultline/vaultline/internal/envlock"
)

func tmpEnvFile(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, ".env")
}

func TestAcquire_CreatesLockFile(t *testing.T) {
	env := tmpEnvFile(t)
	l := envlock.New(env)

	if err := l.Acquire(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer l.Release()

	if _, err := os.Stat(l.Path()); err != nil {
		t.Errorf("lock file not created: %v", err)
	}
}

func TestAcquire_FailsWhenAlreadyLocked(t *testing.T) {
	env := tmpEnvFile(t)
	l1 := envlock.New(env)
	l2 := envlock.New(env)

	if err := l1.Acquire(); err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}
	defer l1.Release()

	err := l2.Acquire()
	if err == nil {
		defer l2.Release()
		t.Fatal("expected error on second acquire, got nil")
	}
	if !isLocked(err) {
		t.Errorf("expected ErrLocked, got %v", err)
	}
}

func TestRelease_RemovesLockFile(t *testing.T) {
	env := tmpEnvFile(t)
	l := envlock.New(env)

	_ = l.Acquire()
	if err := l.Release(); err != nil {
		t.Fatalf("release failed: %v", err)
	}
	if _, err := os.Stat(l.Path()); !os.IsNotExist(err) {
		t.Errorf("lock file still exists after release")
	}
}

func TestRelease_NoopWhenNotLocked(t *testing.T) {
	env := tmpEnvFile(t)
	l := envlock.New(env)

	if err := l.Release(); err != nil {
		t.Errorf("expected no error releasing non-existent lock, got %v", err)
	}
}

func TestAcquire_StaleLockIsOverwritten(t *testing.T) {
	env := tmpEnvFile(t)
	l := envlock.New(env)

	// Write a stale lock with a non-existent PID
	_ = os.WriteFile(l.Path(), []byte("999999999\n2000-01-01T00:00:00Z\n"), 0o600)

	if err := l.Acquire(); err != nil {
		t.Fatalf("expected stale lock to be cleared, got %v", err)
	}
	defer l.Release()
}

func isLocked(err error) bool {
	return err != nil && err.Error() != "" &&
		(err == envlock.ErrLocked || containsStr(err.Error(), "locked"))
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 &&
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
