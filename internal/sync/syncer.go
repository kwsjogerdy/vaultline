package sync

import (
	"fmt"

	"github.com/user/vaultline/internal/audit"
	"github.com/user/vaultline/internal/diff"
	"github.com/user/vaultline/internal/env"
	"github.com/user/vaultline/internal/filter"
	"github.com/user/vaultline/internal/prompt"
	"github.com/user/vaultline/internal/vault"
)

// Syncer orchestrates fetching secrets and writing them to a .env file.
type Syncer struct {
	client   *vault.Client
	writer   *env.Writer
	filter   *filter.Filter
	auditor  *audit.Logger
	prompter *prompt.Prompter
	confirm  bool
}

// New creates a new Syncer.
func New(client *vault.Client, writer *env.Writer, f *filter.Filter, auditor *audit.Logger, prompter *prompt.Prompter, confirm bool) *Syncer {
	return &Syncer{
		client:   client,
		writer:   writer,
		filter:   f,
		auditor:  auditor,
		prompter: prompter,
		confirm:  confirm,
	}
}

// Run fetches secrets from Vault, applies filters, optionally diffs and confirms, then writes.
func (s *Syncer) Run(path, outFile string, existing map[string]string) error {
	secrets, err := s.client.GetSecrets(path)
	if err != nil {
		s.auditor.LogSync(path, outFile, err)
		return fmt.Errorf("fetching secrets: %w", err)
	}

	filtered := s.filter.Apply(secrets)

	if s.confirm {
		changes := diff.Compare(existing, filtered)
		if !changes.HasChanges() {
			fmt.Println("No changes detected.")
			return nil
		}
		ok, err := s.prompter.ConfirmDiff(len(changes.Added), len(changes.Removed), len(changes.Changed))
		if err != nil {
			return fmt.Errorf("prompt: %w", err)
		}
		if !ok {
			fmt.Println("Aborted.")
			return nil
		}
	}

	if err := s.writer.Write(outFile, filtered); err != nil {
		s.auditor.LogSync(path, outFile, err)
		return fmt.Errorf("writing env file: %w", err)
	}

	s.auditor.LogSync(path, outFile, nil)
	return nil
}
