package backup

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Backup creates a timestamped backup of an existing .env file before overwriting it.
// Returns the backup path, or an empty string if the source file did not exist.
func Backup(envPath, backupDir string) (string, error) {
	src, err := os.Open(envPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("backup: open source: %w", err)
	}
	defer src.Close()

	if err := os.MkdirAll(backupDir, 0700); err != nil {
		return "", fmt.Errorf("backup: create dir: %w", err)
	}

	base := filepath.Base(envPath)
	timestamp := time.Now().UTC().Format("20060102T150405Z")
	backupName := fmt.Sprintf("%s.%s.bak", base, timestamp)
	backupPath := filepath.Join(backupDir, backupName)

	dst, err := os.OpenFile(backupPath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0600)
	if err != nil {
		return "", fmt.Errorf("backup: create backup file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("backup: copy: %w", err)
	}

	return backupPath, nil
}

// PruneBackups removes backup files in backupDir older than maxAge.
func PruneBackups(backupDir string, maxAge time.Duration) (int, error) {
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, fmt.Errorf("prune: read dir: %w", err)
	}

	cutoff := time.Now().UTC().Add(-maxAge)
	removed := 0

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		if info.ModTime().Before(cutoff) {
			_ = os.Remove(filepath.Join(backupDir, e.Name()))
			removed++
		}
	}
	return removed, nil
}
