package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultline/internal/config"
	"vaultline/internal/env"
	"vaultline/internal/envdrift"
	"vaultline/internal/vault"
)

var driftEnvFile string

func init() {
	driftCmd := &cobra.Command{
		Use:   "drift",
		Short: "Detect drift between local .env file and Vault secrets",
		RunE:  runDrift,
	}
	driftCmd.Flags().StringVarP(&driftEnvFile, "env-file", "e", ".env", "Path to local .env file")
	rootCmd.AddCommand(driftCmd)
}

func runDrift(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	vaultSecrets, err := client.GetSecrets(cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("vault secrets: %w", err)
	}

	localSecrets, err := env.ReadFile(driftEnvFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read env file: %w", err)
	}
	if localSecrets == nil {
		localSecrets = map[string]string{}
	}

	report := envdrift.Detect(localSecrets, vaultSecrets)

	if !report.HasDrift() {
		fmt.Println("✓ No drift detected.")
		return nil
	}

	fmt.Printf("⚠ Drift detected: %s\n\n", report.Summary())
	for _, e := range report.Entries {
		if e.Status != envdrift.StatusMatch {
			fmt.Printf("  [%-8s] %s\n", e.Status, e.Key)
		}
	}
	return nil
}
