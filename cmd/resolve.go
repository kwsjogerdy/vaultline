package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"vaultline/internal/envresolve"
)

var resolveInput string

func init() {
	resolveCmd := &cobra.Command{
		Use:   "resolve",
		Short: "Interpolate ${KEY} references within a JSON secret map",
		RunE:  runResolve,
	}
	resolveCmd.Flags().StringVarP(&resolveInput, "input", "i", "", "JSON file containing secrets map (required)")
	_ = resolveCmd.MarkFlagRequired("input")
	rootCmd.AddCommand(resolveCmd)
}

func runResolve(cmd *cobra.Command, _ []string) error {
	data, err := os.ReadFile(resolveInput)
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	var secrets map[string]string
	if err := json.Unmarshal(data, &secrets); err != nil {
		return fmt.Errorf("parsing JSON: %w", err)
	}

	r := envresolve.New()
	resolved, err := r.Resolve(secrets)
	if err != nil {
		return fmt.Errorf("resolve: %w", err)
	}

	keys := make([]string, 0, len(resolved))
	for k := range resolved {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Fprintln(cmd.OutOrStdout(), "Resolved secrets:")
	for _, k := range keys {
		fmt.Fprintf(cmd.OutOrStdout(), "  %s=%s\n", k, resolved[k])
	}
	return nil
}
