package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"vaultline/internal/envflatten"
)

var (
	flattenSep    string
	flattenPrefix string
	flattenFile   string
)

func init() {
	flattenCmd := &cobra.Command{
		Use:   "flatten",
		Short: "Flatten a nested JSON map into dot-notation env keys",
		RunE:  runFlatten,
	}
	flattenCmd.Flags().StringVar(&flattenSep, "sep", "__", "Key separator")
	flattenCmd.Flags().StringVar(&flattenPrefix, "prefix", "", "Optional key prefix")
	flattenCmd.Flags().StringVar(&flattenFile, "file", "", "Path to JSON file (default: stdin)")
	rootCmd.AddCommand(flattenCmd)
}

func runFlatten(cmd *cobra.Command, args []string) error {
	var raw []byte
	var err error

	if flattenFile != "" {
		raw, err = os.ReadFile(flattenFile)
	} else {
		raw, err = os.ReadFile("/dev/stdin")
	}
	if err != nil {
		return fmt.Errorf("read input: %w", err)
	}

	var nested map[string]any
	if err := json.Unmarshal(raw, &nested); err != nil {
		return fmt.Errorf("parse JSON: %w", err)
	}

	f := envflatten.New(flattenSep)
	if flattenPrefix != "" {
		f = f.WithPrefix(flattenPrefix)
	}

	out, err := f.Flatten(nested)
	if err != nil {
		return fmt.Errorf("flatten: %w", err)
	}

	keys := make([]string, 0, len(out))
	for k := range out {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, out[k])
	}
	return nil
}
