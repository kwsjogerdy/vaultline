// Package envclip provides functionality to copy individual secret values
// to the system clipboard with optional auto-clear after a timeout.
package envclip

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// Clipper copies secret values to the system clipboard.
type Clipper struct {
	clearAfter time.Duration
}

// New returns a Clipper that clears the clipboard after clearAfter.
// Pass 0 to disable auto-clear.
func New(clearAfter time.Duration) *Clipper {
	return &Clipper{clearAfter: clearAfter}
}

// Copy writes the value for key to the system clipboard.
func (c *Clipper) Copy(secrets map[string]string, key string) error {
	val, ok := secrets[key]
	if !ok {
		return fmt.Errorf("key %q not found", key)
	}
	if err := writeClipboard(val); err != nil {
		return fmt.Errorf("clipboard write: %w", err)
	}
	if c.clearAfter > 0 {
		go func() {
			time.Sleep(c.clearAfter)
			_ = writeClipboard("")
		}()
	}
	return nil
}

func writeClipboard(value string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		cmd = exec.Command("xclip", "-selection", "clipboard")
	case "windows":
		cmd = exec.Command("clip")
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	cmd.Stdin = strings.NewReader(value)
	return cmd.Run()
}
