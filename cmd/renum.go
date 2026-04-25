package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/env"
	"github.com/vaultline/vaultline/internal/envrenum"
)

var (
	renumEnvFile  string
	renumRulesFile string
)

func init() {
	renumCmd := &cobra.Command{
		Use:   "renum",
		Short: "Validate secret values against allowed enum sets",
		RunE:  runRenum,
	}
	renumCmd.Flags().StringVarP(&renumEnvFile, "file", "f", ".env", "Path to .env file")
	renumCmd.Flags().StringVarP(&renumRulesFile, "rules", "r", "", "Path to JSON rules file (required)")
	_ = renumCmd.MarkFlagRequired("rules")
	rootCmd.AddCommand(renumCmd)
}

// rulesFile is the JSON structure for enum rules.
type enumRulesFile struct {
	Rules []struct {
		Key     string   `json:"key"`
		Allowed []string `json:"allowed"`
	} `json:"rules"`
}

func runRenum(cmd *cobra.Command, _ []string) error {
	secrets, err := env.ReadFile(renumEnvFile)
	if err != nil {
		return fmt.Errorf("reading env file: %w", err)
	}

	data, err := os.ReadFile(renumRulesFile)
	if err != nil {
		return fmt.Errorf("reading rules file: %w", err)
	}

	var rf enumRulesFile
	if err := json.Unmarshal(data, &rf); err != nil {
		return fmt.Errorf("parsing rules file: %w", err)
	}

	rules := make([]envrenum.Rule, 0, len(rf.Rules))
	for _, r := range rf.Rules {
		rules = append(rules, envrenum.Rule{Key: r.Key, Allowed: r.Allowed})
	}

	v, err := envrenum.New(rules)
	if err != nil {
		return fmt.Errorf("building validator: %w", err)
	}

	violations := v.Validate(secrets)
	if len(violations) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "✔ all values satisfy enum constraints")
		return nil
	}

	for _, viol := range violations {
		fmt.Fprintln(cmd.ErrOrStderr(), "✘ "+viol.Error())
	}
	return fmt.Errorf("%d enum violation(s) found", len(violations))
}
