package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/env"
	"github.com/vaultline/vaultline/internal/envindex"
)

var (
	indexEnvFile  string
	indexLabelMap string
)

func init() {
	indexCmd := &cobra.Command{
		Use:   "index <query>",
		Short: "Search secrets by key, value, or label substring",
		Args:  cobra.ExactArgs(1),
		RunE:  runIndex,
	}
	indexCmd.Flags().StringVar(&indexEnvFile, "env-file", ".env", "Path to the .env file to index")
	RootCmd.AddCommand(indexCmd)
}

func runIndex(cmd *cobra.Command, args []string) error {
	query := args[0]

	secrets, err := env.ReadFile(indexEnvFile)
	if err != nil {
		return fmt.Errorf("reading env file: %w", err)
	}

	idx := envindex.New(secrets, nil)
	results := idx.Search(query)

	if len(results) == 0 {
		fmt.Fprintf(os.Stderr, "no secrets matched %q\n", query)
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KEY\tMATCHED ON\tVALUE")
	fmt.Fprintln(w, "---\t----------\t-----")
	for _, r := range results {
		val := r.Value
		if len(val) > 40 {
			val = val[:37] + "..."
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", r.Key, r.MatchedOn, val)
	}
	return w.Flush()
}
