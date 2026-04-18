package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/envscope"
)

var scopeCmd = &cobra.Command{
	Use:   "scope",
	Short: "Filter secrets by environment scope",
	RunE:  runScope,
}

var (
	scopeFlag    string
	scopesFile   string
	scopeSecretsFile string
)

func init() {
	scopeCmd.Flags().StringVar(&scopeFlag, "scope", "", "Scope name to apply (required)")
	scopeCmd.Flags().StringVar(&scopesFile, "scopes", "", "JSON file defining scopes")
	scopeCmd.Flags().StringVar(&scopeSecretsFile, "secrets", "", "JSON file of key=value secrets to filter")
	_ = scopeCmd.MarkFlagRequired("scope")
	_ = scopeCmd.MarkFlagRequired("scopes")
	_ = scopeCmd.MarkFlagRequired("secrets")
	rootCmd.AddCommand(scopeCmd)
}

func runScope(cmd *cobra.Command, _ []string) error {
	raw, err := os.ReadFile(scopesFile)
	if err != nil {
		return fmt.Errorf("reading scopes file: %w", err)
	}
	var defs []envscope.Scope
	if err := json.Unmarshal(raw, &defs); err != nil {
		return fmt.Errorf("parsing scopes: %w", err)
	}

	sraw, err := os.ReadFile(scopeSecretsFile)
	if err != nil {
		return fmt.Errorf("reading secrets file: %w", err)
	}
	var secrets map[string]string
	if err := json.Unmarshal(sraw, &secrets); err != nil {
		return fmt.Errorf("parsing secrets: %w", err)
	}

	mgr := envscope.New(defs)
	filtered, err := mgr.Apply(scopeFlag, secrets)
	if err != nil {
		return err
	}

	for k, v := range filtered {
		fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", strings.ToUpper(k), v)
	}
	return nil
}
