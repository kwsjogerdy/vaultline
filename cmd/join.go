package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/envjoin"
)

var (
	joinMode string
	joinSep  string
)

func init() {
	joinCmd := &cobra.Command{
		Use:   "join <left.json> <right.json>",
		Short: "Merge two secret maps using a join strategy",
		Args:  cobra.ExactArgs(2),
		RunE:  runJoin,
	}
	joinCmd.Flags().StringVar(&joinMode, "mode", "concat", "Join mode: concat | left | right")
	joinCmd.Flags().StringVar(&joinSep, "sep", ":", "Separator used in concat mode")
	rootCmd.AddCommand(joinCmd)
}

func runJoin(cmd *cobra.Command, args []string) error {
	left, err := loadJSONSecrets(args[0])
	if err != nil {
		return fmt.Errorf("reading left file: %w", err)
	}
	right, err := loadJSONSecrets(args[1])
	if err != nil {
		return fmt.Errorf("reading right file: %w", err)
	}

	j, err := envjoin.New(envjoin.JoinMode(joinMode), joinSep)
	if err != nil {
		return err
	}

	result, err := j.Join(left, right)
	if err != nil {
		return err
	}

	fmt.Fprint(cmd.OutOrStdout(), envjoin.Summary(left, right, result))
	for _, k := range envjoin.Keys(result) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, result[k])
	}
	return nil
}

func loadJSONSecrets(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	return m, nil
}
