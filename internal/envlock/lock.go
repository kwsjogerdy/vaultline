// Package envlock provides file-based locking to prevent concurrent sync operations
// from corrupting .env files.
package envlock

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var ErrLocked = errors.New("env file is locked by another process")

// Lock represents a file lock for a given env path.
type Lock struct {
	lockPath string
}

// New returns a Lock for the given env file path.
func New(envPath string) *Lock {
	return &Lock{
		lockPath: envPath + ".lock",
	}
}

// Acquire attempts to create a lock file. Returns ErrLocked if already locked.
func (l *Lock) Acquire() error {
	if _, err := os.Stat(l.lockPath); err == nil {
		data, readErr := os.ReadFile(l.lockPath)
		if readErr == nil {
			parts := strings.SplitN(strings.TrimSpace(string(data)), "\n", 2)
			if len(parts) == 2 {
				pid, _ := strconv.Atoi(parts[0])
				if proc, err := os.FindProcess(pid); err == nil {
					if err := proc.Signal(os.Signal(nil)); err == nil {
						return fmt.Errorf("%w: held by pid %d", ErrLocked, pid)
					}
				}
				// stale lock — remove it
				_ = os.Remove(l.lockPath)
			}
		}
	}

	content := fmt.Sprintf("%d\n%s\n", os.Getpid(), time.Now().Format(time.RFC3339))
	if err := os.MkdirAll(filepath.Dir(l.lockPath), 0o755); err != nil {
		return fmt.Errorf("creating lock dir: %w", err)
	}
	if err := os.WriteFile(l.lockPath, []byte(content), 0o600); err != nil {
		return fmt.Errorf("writing lock file: %w", err)
	}
	return nil
}

// Release removes the lock file.
func (l *Lock) Release() error {
	if err := os.Remove(l.lockPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("releasing lock: %w", err)
	}
	return nil
}

// Path returns the lock file path.
func (l *Lock) Path() string {
	return l.lockPath
}
