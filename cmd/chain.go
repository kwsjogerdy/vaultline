package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/envchain"
	"github.com/vaultline/vaultline/internal/env"
)

var chainSources []string
var chainOut string
var chainShowOrigins bool

func init() {
	chainCmd := &cobra.Command{
		Use:   "chain",
		Short: "Merge secrets from multiple .env sources with priority ordering",
		RunE:  runChain,
	}
	chainCmd.Flags().StringArrayVarP(&chainSources, "source", "s", nil, "source files in priority order (lowest first), e.g. -s base.env -s override.env")
	chainCmd.Flags().StringVarP(&chainOut, "out", "o", ".env", "output .env file path")
	chainCmd.Flags().BoolVar(&chainShowOrigins, "origins", false, "print key origins to stdout")
	rootCmd.AddCommand(chainCmd)
}

func runChain(cmd *cobra.Command, args []string) error {
	if len(chainSources) == 0 {
		return fmt.Errorf("at least one --source required")
	}
	var sources []envchain.Source
	for i, path := range chainSources {
		secrets, err := env.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}
		sources = append(sources, envchain.Source{Name: path, Priority: i + 1, Secrets: secrets})
	}
	chain := envchain.New(sources)
	resolved := chain.Resolve()
	if chainShowOrigins {
		origins := make(map[string]string, len(resolved))
		for k := range resolved {
			name, _ := chain.Origin(k)
			origins[k] = name
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(origins)
	}
	w := env.NewWriter(chainOut)
	if err := w.Write(resolved); err != nil {
		return fmt.Errorf("writing output: %w", err)
	}
	fmt.Fprintf(os.Stderr, "chain: wrote %d keys to %s\n", len(resolved), chainOut)
	return nil
}
