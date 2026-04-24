package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/env"
	"github.com/vaultline/vaultline/internal/envgraph"
)

var graphFlags struct {
	envFile string
	format  string
}

func init() {
	graphCmd := &cobra.Command{
		Use:   "graph",
		Short: "Render a dependency graph of secret references",
		Long: `Reads a .env file and analyses ${VAR} interpolation references.

Outputs either a human-readable adjacency list (default) or Graphviz DOT
format suitable for piping into 'dot -Tpng'.`,
		RunE: runGraph,
	}
	graphCmd.Flags().StringVarP(&graphFlags.envFile, "file", "f", ".env", "path to .env file")
	graphCmd.Flags().StringVar(&graphFlags.format, "format", "list", "output format: list or dot")
	rootCmd.AddCommand(graphCmd)
}

func runGraph(cmd *cobra.Command, _ []string) error {
	secrets, err := env.ReadFile(graphFlags.envFile)
	if err != nil {
		return fmt.Errorf("read env file: %w", err)
	}

	g := envgraph.New(secrets)

	switch graphFlags.format {
	case "dot":
		fmt.Fprint(os.Stdout, g.DOT())
	default:
		roots := g.Roots()
		fmt.Fprintf(os.Stdout, "Roots (%d):\n", len(roots))
		for _, r := range roots {
			fmt.Fprintf(os.Stdout, "  %s\n", r)
		}
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Edges:")
		for _, n := range g.Nodes() {
			if len(n.Deps) == 0 {
				continue
			}
			for _, dep := range n.Deps {
				fmt.Fprintf(os.Stdout, "  %s -> %s\n", n.Key, dep)
			}
		}
	}
	return nil
}
