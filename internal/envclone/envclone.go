// Package envclone provides functionality to clone secrets from one
// Vault path to another, optionally remapping keys.
package envclone

import (
	"fmt"
	"strings"
)

// Fetcher retrieves secrets from a source path.
type Fetcher interface {
	GetSecrets(path string) (map[string]string, error)
}

// Writer stores secrets to a destination path.
type Writer interface {
	PutSecrets(path string, secrets map[string]string) error
}

// Options configures a clone operation.
type Options struct {
	// Remap maps source key names to destination key names.
	Remap map[string]string
	// Only copies these keys; empty means copy all.
	Only []string
}

// Cloner copies secrets between Vault paths.
type Cloner struct {
	fetcher Fetcher
	writer  Writer
}

// New creates a new Cloner.
func New(f Fetcher, w Writer) *Cloner {
	return &Cloner{fetcher: f, writer: w}
}

// Clone reads secrets from src and writes them to dst, applying opts.
func (c *Cloner) Clone(src, dst string, opts Options) (int, error) {
	if strings.TrimSpace(src) == "" || strings.TrimSpace(dst) == "" {
		return 0, fmt.Errorf("envclone: src and dst paths must not be empty")
	}

	secrets, err := c.fetcher.GetSecrets(src)
	if err != nil {
		return 0, fmt.Errorf("envclone: fetch from %q: %w", src, err)
	}

	filtered := make(map[string]string)
	onlySet := toSet(opts.Only)

	for k, v := range secrets {
		if len(onlySet) > 0 {
			if _, ok := onlySet[k]; !ok {
				continue
			}
		}
		destKey := k
		if mapped, ok := opts.Remap[k]; ok {
			destKey = mapped
		}
		filtered[destKey] = v
	}

	if err := c.writer.PutSecrets(dst, filtered); err != nil {
		return 0, fmt.Errorf("envclone: write to %q: %w", dst, err)
	}

	return len(filtered), nil
}

func toSet(keys []string) map[string]struct{} {
	s := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		s[k] = struct{}{}
	}
	return s
}
