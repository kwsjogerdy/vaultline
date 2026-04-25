package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/env"
	"github.com/vaultline/vaultline/internal/envstats"
)

var (
	statsFile      string
	statsDelimiter string
)

func init() {
	statsCmd := &cobra.Command{
		Use:   "stats",
		Short: "Print statistics about a .env file",
		RunE:  runStats,
	}
	statsCmd.Flags().StringVarP(&statsFile, "file", "f", ".env", "Path to .env file")
	statsCmd.Flags().StringVarP(&statsDelimiter, "delimiter", "d", "_", "Key prefix delimiter")
	rootCmd.AddCommand(statsCmd)
}

func runStats(cmd *cobra.Command, _ []string) error {
	secrets, err := env.ReadFile(statsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", statsFile)
		}
		return fmt.Errorf("reading env file: %w", err)
	}

	a := envstats.New(statsDelimiter)
	s := a.Compute(secrets)
	fmt.Fprintln(cmd.OutOrStdout(), s.Summary())
	return nil
}
