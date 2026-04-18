// Package multienv supports writing secrets to multiple .env targets.
package multienv

import (
	"fmt"

	"github.com/vaultline/vaultline/internal/env"
)

// Target describes a single output destination with optional key filtering.
type Target struct {
	Path    string
	Keys    []string // if empty, all keys are written
}

// Writer writes secrets to multiple env files.
type Writer struct {
	targets []Target
}

// New creates a new multi-target Writer.
func New(targets []Target) *Writer {
	return &Writer{targets: targets}
}

// WriteAll writes the provided secrets to every configured target,
// filtering keys per target when a key list is specified.
func (w *Writer) WriteAll(secrets map[string]string) error {
	for _, t := range w.targets {
		filtered := filter(secrets, t.Keys)
		wr, err := env.NewWriter(t.Path)
		if err != nil {
			return fmt.Errorf("multienv: open %s: %w", t.Path, err)
		}
		if err := wr.Write(filtered); err != nil {
			return fmt.Errorf("multienv: write %s: %w", t.Path, err)
		}
	}
	return nil
}

// Targets returns the list of targets configured on the Writer.
func (w *Writer) Targets() []Target {
	return w.targets
}

func filter(secrets map[string]string, keys []string) map[string]string {
	if len(keys) == 0 {
		return secrets
	}
	out := make(map[string]string, len(keys))
	for _, k := range keys {
		if v, ok := secrets[k]; ok {
			out[k] = v
		}
	}
	return out
}
