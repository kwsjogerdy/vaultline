package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/envdependency"
)

var dependencyRulesFile string

func init() {
	depCmd := &cobra.Command{
		Use:   "dependency",
		Short: "Check and order secrets based on dependency rules",
	}

	checkCmd := &cobra.Command{
		Use:   "check",
		Short: "Check that all declared dependencies are present in a secrets file",
		RunE:  runDependencyCheck,
	}
	checkCmd.Flags().StringVar(&dependencyRulesFile, "rules", "dependency_rules.json", "Path to dependency rules JSON file")
	checkCmd.Flags().String("env", ".env", "Path to .env file to check")

	orderCmd := &cobra.Command{
		Use:   "order",
		Short: "Print keys in dependency-resolved load order",
		RunE:  runDependencyOrder,
	}
	orderCmd.Flags().StringVar(&dependencyRulesFile, "rules", "dependency_rules.json", "Path to dependency rules JSON file")

	depCmd.AddCommand(checkCmd, orderCmd)
	rootCmd.AddCommand(depCmd)
}

type ruleFile struct {
	Rules []envdependency.Rule `json:"rules"`
}

func loadRules(path string) ([]envdependency.Rule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading rules file: %w", err)
	}
	var rf ruleFile
	if err := json.Unmarshal(data, &rf); err != nil {
		return nil, fmt.Errorf("parsing rules file: %w", err)
	}
	return rf.Rules, nil
}

func runDependencyCheck(cmd *cobra.Command, _ []string) error {
	rules, err := loadRules(dependencyRulesFile)
	if err != nil {
		return err
	}
	envPath, _ := cmd.Flags().GetString("env")
	data, err := os.ReadFile(envPath)
	if err != nil {
		return fmt.Errorf("reading env file: %w", err)
	}
	secrets := make(map[string]string)
	for _, line := range strings.Split(string(data), "\n") {
		parts := strings.SplitN(strings.TrimSpace(line), "=", 2)
		if len(parts) == 2 && parts[0] != "" {
			secrets[parts[0]] = parts[1]
		}
	}
	r, err := envdependency.New(rules)
	if err != nil {
		return err
	}
	results := r.Check(secrets)
	ok := true
	for _, res := range results {
		if !res.Satisfied {
			fmt.Fprintf(cmd.OutOrStdout(), "UNSATISFIED %s — missing: %s\n", res.Key, strings.Join(res.Missing, ", "))
			ok = false
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "OK          %s\n", res.Key)
		}
	}
	if !ok {
		return fmt.Errorf("one or more dependency rules are unsatisfied")
	}
	return nil
}

func runDependencyOrder(cmd *cobra.Command, _ []string) error {
	rules, err := loadRules(dependencyRulesFile)
	if err != nil {
		return err
	}
	order, err := envdependency.Order(rules)
	if err != nil {
		return err
	}
	for _, k := range order {
		fmt.Fprintln(cmd.OutOrStdout(), k)
	}
	return nil
}
