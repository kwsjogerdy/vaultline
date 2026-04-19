package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/env"
	"github.com/vaultline/vaultline/internal/envdiff"
)

var envdiffCmd = &cobra.Command{
	Use:   "envdiff <left> <right>",
	Short: "Compare two .env files side by side",
	Args:  cobra.ExactArgs(2),
	RunE:  runEnvDiff,
}

func init() {
	rootCmd.AddCommand(envdiffCmd)
}

func runEnvDiff(cmd *cobra.Command, args []string) error {
	left, err := env.ReadFile(args[0])
	if err != nil {
		return fmt.Errorf("reading left file: %w", err)
	}
	right, err := env.ReadFile(args[1])
	if err != nil {
		return fmt.Errorf("reading right file: %w", err)
	}

	result := envdiff.Compare(left, right)

	statusIcon := map[string]string{
		"match":      "=",
		"changed":    "~",
		"left_only":  "<",
		"right_only": ">",
	}

	for _, e := range result.Entries {
		icon := statusIcon[e.Status]
		switch e.Status {
		case "match":
			fmt.Fprintf(os.Stdout, "  %s %s\n", icon, e.Key)
		case "changed":
			fmt.Fprintf(os.Stdout, "  %s %s: %q -> %q\n", icon, e.Key, e.Left, e.Right)
		case "left_only":
			fmt.Fprintf(os.Stdout, "  %s %s (left only)\n", icon, e.Key)
		case "right_only":
			fmt.Fprintf(os.Stdout, "  %s %s (right only)\n", icon, e.Key)
		}
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, result.Summary())

	if result.HasDrift() {
		os.Exit(1)
	}
	return nil
}
