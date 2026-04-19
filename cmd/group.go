package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultline/internal/envgroup"
	"vaultline/internal/env"
)

var groupSep string
var groupEnvFile string

func init() {
	groupCmd := &cobra.Command{
		Use:   "group",
		Short: "Display secrets grouped by key prefix",
		RunE:  runGroup,
	}
	groupCmd.Flags().StringVar(&groupSep, "sep", "_", "Delimiter used to split key prefixes")
	groupCmd.Flags().StringVar(&groupEnvFile, "file", ".env", "Path to the .env file to read")
	rootCmd.AddCommand(groupCmd)
}

func runGroup(cmd *cobra.Command, args []string) error {
	secrets, err := env.ReadFile(groupEnvFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", groupEnvFile)
		}
		return fmt.Errorf("reading env file: %w", err)
	}

	if len(secrets) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no secrets found")
		return nil
	}

	g := envgroup.New(groupSep)
	groups := g.Group(secrets)

	for _, gr := range groups {
		fmt.Fprintf(cmd.OutOrStdout(), "\n[%s] — %d key(s)\n", gr.Name, len(gr.Secrets))
		for k := range gr.Secrets {
			fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", k)
		}
	}
	return nil
}
