package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/env"
	"github.com/vaultline/vaultline/internal/envredact"
)

var (
	redactInput       string
	redactPatterns    []string
	redactPlaceholder string
	redactShowKeys    bool
)

func init() {
	redactCmd := &cobra.Command{
		Use:   "redact",
		Short: "Print secrets with sensitive values replaced by a placeholder",
		RunE:  runRedact,
	}
	redactCmd.Flags().StringVarP(&redactInput, "file", "f", ".env", "Path to the .env file to redact")
	redactCmd.Flags().StringSliceVar(&redactPatterns, "pattern", nil, "Additional key patterns to redact (regex)")
	redactCmd.Flags().StringVar(&redactPlaceholder, "placeholder", "[REDACTED]", "Replacement text for sensitive values")
	redactCmd.Flags().BoolVar(&redactShowKeys, "show-keys", false, "Print the list of redacted keys")
	rootCmd.AddCommand(redactCmd)
}

func runRedact(cmd *cobra.Command, _ []string) error {
	secrets, err := env.ReadFile(redactInput)
	if err != nil {
		return fmt.Errorf("read %s: %w", redactInput, err)
	}

	var r *envredact.Redactor
	if len(redactPatterns) > 0 {
		r = envredact.NewWithPatterns(redactPatterns, redactPlaceholder)
	} else {
		r = envredact.New()
	}

	result := r.Apply(secrets)

	keys := make([]string, 0, len(result.Secrets))
	for k := range result.Secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, result.Secrets[k])
	}

	if redactShowKeys && len(result.Redacted) > 0 {
		sort.Strings(result.Redacted)
		fmt.Fprintf(cmd.OutOrStdout(), "\nRedacted keys (%d):\n", len(result.Redacted))
		for _, k := range result.Redacted {
			fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", k)
		}
	}
	return nil
}
