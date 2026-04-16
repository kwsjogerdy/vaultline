package sync

import (
	"fmt"

	"github.com/vaultline/internal/env"
	"github.com/vaultline/internal/vault"
)

// Result holds the outcome of a sync operation.
type Result struct {
	Path    string
	Count   int
	Skipped int
}

// Syncer orchestrates fetching secrets from Vault and writing them to a .env file.
type Syncer struct {
	client *vault.Client
	writer *env.Writer
}

// New creates a new Syncer.
func New(client *vault.Client, writer *env.Writer) *Syncer {
	return &Syncer{client: client, writer: writer}
}

// Sync fetches secrets at the given Vault path and writes them to outPath.
// Keys listed in skip are excluded from the output.
func (s *Syncer) Sync(vaultPath, outPath string, skip []string) (*Result, error) {
	secrets, err := s.client.GetSecrets(vaultPath)
	if err != nil {
		return nil, fmt.Errorf("fetching secrets from %q: %w", vaultPath, err)
	}

	skipSet := make(map[string]struct{}, len(skip))
	for _, k := range skip {
		skipSet[k] = struct{}{}
	}

	filtered := make(map[string]string, len(secrets))
	skipped := 0
	for k, v := range secrets {
		if _, excluded := skipSet[k]; excluded {
			skipped++
			continue
		}
		filtered[k] = v
	}

	if err := s.writer.Write(outPath, filtered); err != nil {
		return nil, fmt.Errorf("writing env file %q: %w", outPath, err)
	}

	return &Result{
		Path:    outPath,
		Count:   len(filtered),
		Skipped: skipped,
	}, nil
}
