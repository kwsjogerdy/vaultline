// Package watch provides file-system watching to trigger re-sync when
// the vaultline config file changes.
package watch

import (
	"context"
	"log"
	"os"
	"time"
)

// SyncFunc is called whenever a change is detected.
type SyncFunc func() error

// Watcher polls a file for modifications and calls fn on change.
type Watcher struct {
	path     string
	interval time.Duration
	fn       SyncFunc
}

// New creates a Watcher for the given file path and poll interval.
func New(path string, interval time.Duration, fn SyncFunc) *Watcher {
	return &Watcher{path: path, interval: interval, fn: fn}
}

// Run starts the polling loop and blocks until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	info, err := os.Stat(w.path)
	if err != nil {
		return err
	}
	lastMod := info.ModTime()

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			info, err := os.Stat(w.path)
			if err != nil {
				log.Printf("watch: stat error: %v", err)
				continue
			}
			if info.ModTime().After(lastMod) {
				lastMod = info.ModTime()
				log.Printf("watch: change detected in %s, triggering sync", w.path)
				if err := w.fn(); err != nil {
					log.Printf("watch: sync error: %v", err)
				}
			}
		}
	}
}
