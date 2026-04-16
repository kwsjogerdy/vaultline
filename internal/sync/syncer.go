package sync

import (
	"fmt"

	"github.com/vaultline/vaultline/internal/env"
	"github.com/vaultline/vaultline/internal/filter"
	"github.com/vaultline/vaultline/internal/vault"
)

// Syncer orchestrates fetching secrets from Vault and writing them to a file.
type Syncer struct {
	client *vault.Client
	writer *env.Writer
	filter *filter.Filter
}

// New creates a Syncer with the provided client, writer, and filter rule.
func New(client *vault.Client, writer *env.Writer, rule filter.Rule) *Syncer {
	return &Syncer{
		client: client,
		writer: writer,
		filter: filter.New(rule),
	}
}

// Run fetches secrets from the given path, applies filters, and writes to the output file.
func (s *Syncer) Run(secretPath string) error {
	secrets, err := s.client.GetSecrets(secretPath)
	if err != nil {
		return fmt.Errorf("fetching secrets: %w", err)
	}

	filtered := s.filter.Apply(secrets)
	if len(filtered) == 0 {
		return fmt.Errorf("no secrets matched the filter rules for path %q", secretPath)
	}

	if err := s.writer.Write(filtered); err != nil {
		return fmt.Errorf("writing secrets: %w", err)
	}

	return nil
}
