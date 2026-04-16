package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vaultline",
	Short: "Sync secrets from Vault to local .env files",
	Long: `vaultline is a lightweight CLI tool for securely syncing
secrets from HashiCorp Vault into local .env files with
diff previups, and audit logging.`,
	SilenceUsage: true,
}

 Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
	intln(os.Stderr, err)
		os.Exit(1)
	}
}
