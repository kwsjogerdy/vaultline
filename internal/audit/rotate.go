package audit

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// RotateOptions controls log rotation behaviour.
type RotateOptions struct {
	// MaxBytes is the maximum file size before rotation (0 = no limit).
	MaxBytes int64
	// Keep is the number of rotated files to retain (0 = keep all).
	Keep int
}

// globFn is a thin wrapper so tests can stub it.
var globFn = filepath.Glob

// Rotate renames logPath to a timestamped backup when it exceeds MaxBytes,
// then prunes old rotated files when Keep > 0.
func Rotate(logPath string, opts RotateOptions) error {
	if opts.MaxBytes <= 0 {
		return nil
	}
	info, err := os.Stat(logPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("audit rotate: stat: %w", err)
	}
	if info.Size() < opts.MaxBytes {
		return nil
	}

	ts := time.Now().UTC().Format("20060102T150405")
	dest := fmt.Sprintf("%s.%s", logPath, ts)
	if err := os.Rename(logPath, dest); err != nil {
		return fmt.Errorf("audit rotate: rename: %w", err)
	}

	if opts.Keep > 0 {
		return pruneOldLogs(logPath, opts.Keep)
	}
	return nil
}

func pruneOldLogs(base string, keep int) error {
	matches, err := globFn(base + ".*")
	if err != nil {
		return fmt.Errorf("audit rotate: glob: %w", err)
	}
	if len(matches) <= keep {
		return nil
	}
	for _, old := range matches[:len(matches)-keep] {
		if err := os.Remove(old); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("audit rotate: remove %s: %w", old, err)
		}
	}
	return nil
}
