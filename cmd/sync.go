package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"vaultline/internal/audit"
	"vaultline/internal/backup"
	"vaultline/internal/config"
	"vaultline/internal/diff"
	"vaultline/internal/env"
	"vaultline/internal/filter"
	"vaultline/internal/prompt"
	"vaultline/internal/sync"
	"vaultline/internal/vault"
)

var (
	flagDryRun  bool
	flagYes     bool
	flagPrefix  string
	flagExclude []string
	flagOutput  string
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync secrets from Vault to a local .env file",
	RunE:  runSync,
}

func init() {
	syncCmd.Flags().BoolVar(&flagDryRun, "dry-run", false, "Preview changes without writing")
	syncCmd.Flags().BoolVarP(&flagYes, "yes", "y", false, "Skip confirmation prompt")
	syncCmd.Flags().StringVar(&flagPrefix, "prefix", "", "Only include keys with this prefix")
	syncCmd.Flags().StringArrayVar(&flagExclude, "exclude", nil, "Keys to exclude")
	syncCmd.Flags().StringVarP(&flagOutput, "output", "o", ".env", "Output .env file path")
	rootCmd.AddCommand(syncCmd)
}

func runSync(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	f := filter.New(flagPrefix, flagExclude)
	w := env.NewWriter(flagOutput)
	auditor := audit.New(cfg.AuditLog)
	syncer := sync.New(client, w, f, auditor)

	secrets, err := syncer.FetchFiltered()
	if err != nil {
		return fmt.Errorf("fetch: %w", err)
	}

	existing, _ := env.ReadFile(flagOutput)
	changes := diff.Compare(existing, secrets)

	if !diff.HasChanges(changes) {
		fmt.Println("No changes detected.")
		return nil
	}

	p := prompt.New()
	if !flagYes && !flagDryRun {
		ok, err := p.ConfirmDiff(changes)
		if err != nil || !ok {
			fmt.Println("Aborted.")
			return nil
		}
	}

	if flagDryRun {
		fmt.Println("Dry run — no changes written.")
		return nil
	}

	if err := backup.Backup(flagOutput, cfg.BackupDir); err != nil {
		fmt.Fprintf(os.Stderr, "warning: backup failed: %v\n", err)
	}

	return syncer.Write(secrets)
}
