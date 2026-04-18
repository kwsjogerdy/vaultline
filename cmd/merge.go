package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/env"
	"github.com/vaultline/vaultline/internal/envmerge"
)

var mergeStrategy string
var mergeOutput string

func init() {
	mergeCmd := &cobra.Command{
		Use:   "merge <file1> <file2> [...]",
		Short: "Merge multiple .env files into one",
		Args:  cobra.MinimumNArgs(2),
		RunE:  runMerge,
	}
	mergeCmd.Flags().StringVarP(&mergeStrategy, "strategy", "s", "first", "Conflict strategy: first, last, error")
	mergeCmd.Flags().StringVarP(&mergeOutput, "output", "o", ".env.merged", "Output file path")
	rootCmd.AddCommand(mergeCmd)
}

func runMerge(cmd *cobra.Command, args []string) error {
	var strategy envmerge.Strategy
	switch mergeStrategy {
	case "last":
		strategy = envmerge.StrategyLast
	case "error":
		strategy = envmerge.StrategyError
	default:
		strategy = envmerge.StrategyFirst
	}

	sources := make([]map[string]string, 0, len(args))
	for _, f := range args {
		m, err := env.ReadFile(f)
		if err != nil {
			return fmt.Errorf("reading %s: %w", f, err)
		}
		sources = append(sources, m)
	}

	merger := envmerge.New(strategy)
	res, err := merger.Merge(sources...)
	if err != nil {
		return err
	}

	if len(res.Conflicts) > 0 {
		fmt.Fprintf(os.Stderr, "conflicts (%d):\n", len(res.Conflicts))
		for _, c := range res.Conflicts {
			b, _ := json.Marshal(c)
			fmt.Fprintf(os.Stderr, "  %s\n", b)
		}
	}

	w, err := env.NewWriter(mergeOutput)
	if err != nil {
		return err
	}
	if err := w.Write(res.Secrets); err != nil {
		return err
	}
	fmt.Printf("merged %d keys → %s\n", len(res.Secrets), mergeOutput)
	return nil
}
