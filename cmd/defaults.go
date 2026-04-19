package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/envdefaults"
	"github.com/vaultline/vaultline/internal/env"
)

var defaultsCmd = &cobra.Command{
	Use:   "defaults",
	Short: "Apply default values to missing keys in a .env file",
	RunE:  runDefaults,
}

func init() {
	defaultsCmd.Flags().StringP("file", "f", ".env", "Target .env file")
	defaultsCmd.Flags().StringP("defaults", "d", "defaults.json", "JSON file of default values")
	defaultsCmd.Flags().Bool("show-missing", false, "Print keys that were missing and filled")
	rootCmd.AddCommand(defaultsCmd)
}

func runDefaults(cmd *cobra.Command, _ []string) error {
	envFile, _ := cmd.Flags().GetString("file")
	defaultsFile, _ := cmd.Flags().GetString("defaults")
	showMissing, _ := cmd.Flags().GetBool("show-missing")

	secrets, err := env.ReadFile(envFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("reading env file: %w", err)
	}
	if secrets == nil {
		secrets = map[string]string{}
	}

	defs, err := envdefaults.LoadFile(defaultsFile)
	if err != nil {
		return fmt.Errorf("loading defaults: %w", err)
	}

	applier, err := envdefaults.New(defs)
	if err != nil {
		return err
	}

	if showMissing {
		missing := applier.Missing(secrets)
		if len(missing) > 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "Filling defaults for: %s\n", strings.Join(missing, ", "))
		}
	}

	result := applier.Apply(secrets)
	w := env.NewWriter(envFile)
	if err := w.Write(result); err != nil {
		return fmt.Errorf("writing env file: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Defaults applied to %s\n", envFile)
	return nil
}
