package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/envimport"
)

var (
	importEnvFile  string
	importVaultPath string
	importSkip     []string
)

func init() {
	importCmd := &cobra.Command{
		Use:   "import",
		Short: "Import secrets from a .env file into Vault",
		RunE:  runImport,
	}
	importCmd.Flags().StringVarP(&importEnvFile, "file", "f", ".env", "Path to the .env file to import")
	importCmd.Flags().StringVarP(&importVaultPath, "path", "p", "", "Vault secret path (e.g. secret/data/myapp)")
	importCmd.Flags().StringSliceVar(&importSkip, "skip", nil, "Keys to skip during import")
	_ = importCmd.MarkFlagRequired("path")
	rootCmd.AddCommand(importCmd)
}

func runImport(cmd *cobra.Command, _ []string) error {
	vaultAddr := os.Getenv("VAULT_ADDR")
	if vaultAddr == "" {
		vaultAddr = "http://127.0.0.1:8200"
	}
	vaultToken := os.Getenv("VAULT_TOKEN")
	if vaultToken == "" {
		return fmt.Errorf("VAULT_TOKEN is not set")
	}

	im := envimport.New(vaultAddr, vaultToken)
	res, err := im.FromFile(importEnvFile, importVaultPath, importSkip)
	if err != nil {
		return fmt.Errorf("import failed: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Imported: %d  Skipped: %d  Errors: %d\n",
		len(res.Imported), len(res.Skipped), len(res.Errors))

	for _, k := range res.Imported {
		fmt.Fprintf(cmd.OutOrStdout(), "  ✓ %s\n", k)
	}
	for _, k := range res.Skipped {
		fmt.Fprintf(cmd.OutOrStdout(), "  - %s (skipped)\n", k)
	}
	for _, e := range res.Errors {
		fmt.Fprintf(cmd.ErrOrStderr(), "  ✗ %s\n", e)
	}
	if len(res.Errors) > 0 {
		return fmt.Errorf("%d secret(s) failed to import", len(res.Errors))
	}
	return nil
}
