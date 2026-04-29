package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/envpivot"
)

var pivotCollision string

func init() {
	pivotCmd := &cobra.Command{
		Use:   "pivot",
		Short: "Swap keys and values in a secrets JSON file",
		RunE:  runPivot,
	}
	pivotCmd.Flags().StringVar(&pivotCollision, "collision", "error",
		"Strategy for duplicate values: error | skip | suffix")
	rootCmd.AddCommand(pivotCmd)
}

func runPivot(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: vaultline pivot <secrets.json>")
	}

	data, err := os.ReadFile(args[0])
	if err != nil {
		return fmt.Errorf("pivot: read file: %w", err)
	}

	var secrets map[string]string
	if err := json.Unmarshal(data, &secrets); err != nil {
		return fmt.Errorf("pivot: parse JSON: %w", err)
	}

	p, err := envpivot.New(envpivot.Options{
		Collision: envpivot.CollisionStrategy(pivotCollision),
	})
	if err != nil {
		return fmt.Errorf("pivot: %w", err)
	}

	result, err := p.Apply(secrets)
	if err != nil {
		return fmt.Errorf("pivot: %w", err)
	}

	skipped := p.Skipped(secrets)
	if len(skipped) > 0 {
		fmt.Fprintf(cmd.ErrOrStderr(), "skipped keys (duplicate values): %v\n", skipped)
	}

	keys := make([]string, 0, len(result))
	for k := range result {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, result[k])
	}
	return nil
}
