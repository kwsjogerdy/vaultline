package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/env"
	"github.com/vaultline/vaultline/internal/envmask"
)

var maskCmd = &cobra.Command{
	Use:   "mask",
	Short: "Display secrets with sensitive values masked",
	RunE:  runMask,
}

func init() {
	maskCmd.Flags().String("file", ".env", "Path to the .env file")
	maskCmd.Flags().StringArray("pattern", []string{"password", "secret", "token", "key", "apikey"}, "Key patterns to mask (regex)")
	maskCmd.Flags().Int("reveal", 0, "Number of trailing characters to reveal (0 = full redact)")
	maskCmd.Flags().Bool("json", false, "Output as JSON")
	rootCmd.AddCommand(maskCmd)
}

func runMask(cmd *cobra.Command, _ []string) error {
	filePath, _ := cmd.Flags().GetString("file")
	patterns, _ := cmd.Flags().GetStringArray("pattern")
	reveal, _ := cmd.Flags().GetInt("reveal")
	jsonOut, _ := cmd.Flags().GetBool("json")

	secrets, err := env.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("reading env file: %w", err)
	}

	rules := make([]envmask.Rule, len(patterns))
	for i, p := range patterns {
		rules[i] = envmask.Rule{Pattern: p, Reveal: reveal}
	}

	m, err := envmask.New(rules)
	if err != nil {
		return fmt.Errorf("building masker: %w", err)
	}

	masked := m.Apply(secrets)

	if jsonOut {
		return json.NewEncoder(os.Stdout).Encode(masked)
	}

	keys := make([]string, 0, len(masked))
	for k := range masked {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%s=%s\n", k, masked[k])
	}
	return nil
}
