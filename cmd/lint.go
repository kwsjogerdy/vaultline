package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/env"
	"github.com/vaultline/vaultline/internal/envlint"
)

var lintAllowEmpty bool

func init() {
	lintCmd := &cobra.Command{
		Use:   "lint [env-file]",
		Short: "Lint an .env file for common issues",
		Args:  cobra.ExactArgs(1),
		RunE:  runLint,
	}
	lintCmd.Flags().BoolVar(&lintAllowEmpty, "allow-empty", false, "allow empty values without warning")
	rootCmd.AddCommand(lintCmd)
}

func runLint(cmd *cobra.Command, args []string) error {
	path := args[0]

	secrets, err := env.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading env file: %w", err)
	}

	l := envlint.New(lintAllowEmpty)
	findings := l.Lint(secrets)

	if len(findings) == 0 {
		fmt.Println("✓ No issues found.")
		return nil
	}

	for _, f := range findings {
		fmt.Println(f.String())
	}

	if envlint.HasErrors(findings) {
		fmt.Fprintf(os.Stderr, "\n%d issue(s) found — errors present.\n", len(findings))
		os.Exit(1)
	}

	fmt.Printf("\n%d warning(s) found.\n", len(findings))
	return nil
}
