package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/envshuffle"
)

func init() {
	var (
		inputFile string
		seed      int64
		prefix    string
		reverse   bool
		outputJSON bool
	)

	cmd := &cobra.Command{
		Use:   "shuffle",
		Short: "Randomly reorder secrets from a JSON file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runShuffle(inputFile, seed, prefix, reverse, outputJSON)
		},
	}

	cmd.Flags().StringVarP(&inputFile, "file", "f", "", "JSON file of secrets (required)")
	cmd.Flags().Int64Var(&seed, "seed", 0, "Random seed for deterministic shuffle")
	cmd.Flags().StringVar(&prefix, "prefix", "", "Only shuffle keys matching this prefix")
	cmd.Flags().BoolVar(&reverse, "reverse", false, "Return keys in reverse-sorted order")
	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output as JSON array")
	_ = cmd.MarkFlagRequired("file")

	rootCmd.AddCommand(cmd)
}

func runShuffle(inputFile string, seed int64, prefix string, reverse bool, outputJSON bool) error {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("shuffle: read file: %w", err)
	}

	var secrets map[string]string
	if err := json.Unmarshal(data, &secrets); err != nil {
		return fmt.Errorf("shuffle: parse JSON: %w", err)
	}

	opts := []envshuffle.Option{envshuffle.WithSeed(seed)}
	if prefix != "" {
		opts = append(opts, envshuffle.WithPrefix(prefix))
	}
	if reverse {
		opts = append(opts, envshuffle.WithReverse())
	}

	s := envshuffle.New(opts...)
	entries, err := s.Apply(secrets)
	if err != nil {
		return fmt.Errorf("shuffle: %w", err)
	}

	if outputJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(entries)
	}

	for _, e := range entries {
		fmt.Printf("%s=%s\n", e.Key, e.Value)
	}
	return nil
}
