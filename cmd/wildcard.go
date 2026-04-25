package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/env"
	"github.com/vaultline/vaultline/internal/envwildcard"
)

var (
	wildcardInclude []string
	wildcardExclude []string
	wildcardFile    string
)

func init() {
	wildcardCmd := &cobra.Command{
		Use:   "wildcard",
		Short: "Filter secrets using glob-style wildcard patterns",
		Long: `Filter secrets from a .env file using glob patterns.

Examples:
  vaultline wildcard --file .env --include "DB_*" --exclude "*_PASSWORD"
  vaultline wildcard --file .env --include "APP_*" --include "REDIS_*"`,
		RunE: runWildcard,
	}
	wildcardCmd.Flags().StringSliceVar(&wildcardInclude, "include", nil, "Glob patterns to include (repeatable)")
	wildcardCmd.Flags().StringSliceVar(&wildcardExclude, "exclude", nil, "Glob patterns to exclude (repeatable)")
	wildcardCmd.Flags().StringVar(&wildcardFile, "file", ".env", "Source .env file")
	rootCmd.AddCommand(wildcardCmd)
}

func runWildcard(cmd *cobra.Command, _ []string) error {
	secrets, err := env.ReadFile(wildcardFile)
	if err != nil {
		return fmt.Errorf("read %s: %w", wildcardFile, err)
	}

	m, err := envwildcard.New(wildcardInclude, wildcardExclude)
	if err != nil {
		return fmt.Errorf("wildcard: %w", err)
	}

	matched, err := m.Matched(secrets)
	if err != nil {
		return fmt.Errorf("wildcard match: %w", err)
	}

	if len(matched) == 0 {
		fmt.Fprintln(os.Stderr, "no keys matched the given patterns")
		return nil
	}

	out, _ := m.Apply(secrets)
	w := cmd.OutOrStdout()
	for _, k := range matched {
		fmt.Fprintf(w, "%s=%s\n", k, out[k])
	}
	fmt.Fprintf(os.Stderr, "\n%d key(s) matched (include: [%s], exclude: [%s])\n",
		len(matched),
		strings.Join(wildcardInclude, ", "),
		strings.Join(wildcardExclude, ", "),
	)
	return nil
}
