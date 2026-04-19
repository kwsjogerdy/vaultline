package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/env"
	"github.com/vaultline/vaultline/internal/envreport"
)

var (
	reportFormat string
	reportInput  string
)

func init() {
	reportCmd := &cobra.Command{
		Use:   "report",
		Short: "Generate a summary report of a .env file",
		RunE:  runReport,
	}
	reportCmd.Flags().StringVarP(&reportInput, "input", "i", ".env", "Path to .env file")
	reportCmd.Flags().StringVarP(&reportFormat, "format", "f", "text", "Output format: text or markdown")
	rootCmd.AddCommand(reportCmd)
}

func runReport(cmd *cobra.Command, args []string) error {
	secrets, err := env.ReadFile(reportInput)
	if err != nil {
		return fmt.Errorf("reading %s: %w", reportInput, err)
	}

	fmt := envreport.Format(reportFormat)
	if fmt != envreport.FormatText && fmt != envreport.FormatMarkdown {
		return fmt.Errorf("unsupported format %q: use text or markdown", reportFormat)
	}

	rep := envreport.Build(secrets)
	r := envreport.New(os.Stdout, fmt)
	return r.Write(rep)
}
