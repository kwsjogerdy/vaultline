package rollback

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Manager handles rollback of .env files from backups.
type Manager struct {
	backupDir string
}

// New creates a new rollback Manager.
func New(backupDir string) *Manager {
	return &Manager{backupDir: backupDir}
}

// ListBackups returns available backup files for the given env file, newest first.
func (m *Manager) ListBackups(envFile string) ([]string, error) {
	base := filepath.Base(envFile)
	entries, err := os.ReadDir(m.backupDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read backup dir: %w", err)
	}

	var matches []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasPrefix(e.Name(), base+".") {
			matches = append(matches, filepath.Join(m.backupDir, e.Name()))
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(matches)))
	return matches, nil
}

// Rollback restores the most recent backup of envFile, or a specific backup if backupPath is set.
func (m *Manager) Rollback(envFile, backupPath string) (string, error) {
	if backupPath == "" {
		backups, err := m.ListBackups(envFile)
		if err != nil {
			return "", err
		}
		if len(backups) == 0 {
			return "", fmt.Errorf("no backups found for %s", envFile)
		}
		backupPath = backups[0]
	}

	if err := copyFile(backupPath, envFile); err != nil {
		return "", fmt.Errorf("restore backup: %w", err)
	}
	return backupPath, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
