package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vaultline/vaultline/internal/audit"
	"github.com/vaultline/vaultline/internal/prompt"
	"github.com/vaultline/vaultline/internal/rollback"
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Restore the previous version of your .env file from a backup",
	RunE:  runRollback,
}

var rollbackBackupPath string

func init() {
	rootCmd.AddCommand(rollbackCmd)
	rollbackCmd.Flags().StringVar(&rollbackBackupPath, "backup", "", "specific backup file to restore")
}

func runRollback(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	m := rollback.New(cfg.BackupDir)

	if rollbackBackupPath == "" {
		backups, err := m.ListBackups(cfg.EnvFile)
		if err != nil {
			return fmt.Errorf("list backups: %w", err)
		}
		if len(backups) == 0 {
			fmt.Fprintln(os.Stderr, "No backups available.")
			return nil
		}
		fmt.Println("Available backups (newest first):")
		for i, b := range backups {
			fmt.Printf("  [%d] %s\n", i+1, b)
		}
	}

	p := prompt.New()
	ok, err := p.Confirm(fmt.Sprintf("Restore %s from backup?", cfg.EnvFile))
	if err != nil || !ok {
		fmt.Println("Rollback cancelled.")
		return nil
	}

	used, err := m.Rollback(cfg.EnvFile, rollbackBackupPath)
	if err != nil {
		return fmt.Errorf("rollback failed: %w", err)
	}

	fmt.Printf("Restored %s from %s\n", cfg.EnvFile, used)

	logger, _ := audit.New(cfg.AuditLog)
	logger.Log(audit.Entry{Event: "rollback", Detail: used})

	return nil
}
