package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/config"
	"github.com/vaultline/vaultline/internal/multienv"
	"github.com/vaultline/vaultline/internal/vault"
)

var multienvTargetsFile string

func init() {
	multienvCmd := &cobra.Command{
		Use:   "multienv",
		Short: "Sync secrets to multiple .env files",
		RunE:  runMultienv,
	}
	multienvCmd.Flags().StringVar(&multienvTargetsFile, "targets", "targets.json", "JSON file describing output targets")
	rootCmd.AddCommand(multienvCmd)
}

func runMultienv(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	raw, err := os.ReadFile(multienvTargetsFile)
	if err != nil {
		return fmt.Errorf("read targets file: %w", err)
	}

	var targets []multienv.Target
	if err := json.Unmarshal(raw, &targets); err != nil {
		return fmt.Errorf("parse targets: %w", err)
	}

	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	secrets, err := client.GetSecrets(cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("get secrets: %w", err)
	}

	w := multienv.New(targets)
	if err := w.WriteAll(secrets); err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "wrote secrets to %d target(s)\n", len(targets))
	return nil
}
