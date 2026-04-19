package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/envrewrite"
	"github.com/vaultline/vaultline/internal/env"
)

var rewriteRulesFile string
var rewriteEnvFile string
var rewriteOutput string

func init() {
	rewriteCmd := &cobra.Command{
		Use:   "rewrite",
		Short: "Apply find-and-replace rules to keys/values in a .env file",
		RunE:  runRewrite,
	}
	rewriteCmd.Flags().StringVar(&rewriteRulesFile, "rules", "", "JSON file containing rewrite rules (required)")
	rewriteCmd.Flags().StringVar(&rewriteEnvFile, "env", ".env", "Source .env file")
	rewriteCmd.Flags().StringVar(&rewriteOutput, "out", "", "Output file (defaults to source file)")
	_ = rewriteCmd.MarkFlagRequired("rules")
	rootCmd.AddCommand(rewriteCmd)
}

func runRewrite(cmd *cobra.Command, args []string) error {
	rulesData, err := os.ReadFile(rewriteRulesFile)
	if err != nil {
		return fmt.Errorf("reading rules file: %w", err)
	}
	var rules []envrewrite.Rule
	if err := json.Unmarshal(rulesData, &rules); err != nil {
		return fmt.Errorf("parsing rules: %w", err)
	}

	secrets, err := env.ReadFile(rewriteEnvFile)
	if err != nil {
		return fmt.Errorf("reading env file: %w", err)
	}

	r := envrewrite.New(rules)
	result, err := r.Apply(secrets)
	if err != nil {
		return err
	}

	out := rewriteOutput
	if out == "" {
		out = rewriteEnvFile
	}
	w := env.NewWriter(out)
	if err := w.Write(result); err != nil {
		return fmt.Errorf("writing output: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Rewrote %d keys → %s\n", len(result), out)
	return nil
}
