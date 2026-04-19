package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/env"
	"github.com/vaultline/vaultline/internal/envtrim"
)

var trimPatterns []string

func init() {
	trimCmd := &cobra.Command{
		Use:   "trim",
		Short: "Remove placeholder/blank secrets from a .env file",
		RunE:  runTrim,
	}
	trimCmd.Flags().StringVarP(&envFile, "file", "f", ".env", "Path to .env file")
	trimCmd.Flags().StringSliceVar(&trimPatterns, "pattern", nil, "Extra placeholder patterns to remove (comma-separated)")
	trimCmd.Flags().BoolP("dry-run", "n", false, "Print removed keys without modifying the file")
	rootCmd.AddCommand(trimCmd)
}

func runTrim(cmd *cobra.Command, _ []string) error {
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	secrets, err := env.ReadFile(envFile)
	if err != nil {
		return fmt.Errorf("reading env file: %w", err)
	}

	var tr *envtrim.Trimmer
	if len(trimPatterns) > 0 {
		all := append(envtrim.DefaultPatterns, trimPatterns...)
		tr = envtrim.NewWithPatterns(all)
	} else {
		tr = envtrim.New()
	}

	removed := tr.Removed(secrets)
	sort.Strings(removed)

	if len(removed) == 0 {
		fmt.Println("No placeholder secrets found.")
		return nil
	}

	for _, k := range removed {
		fmt.Printf("  - %s\n", k)
	}

	if dryRun {
		fmt.Printf("[dry-run] %d key(s) would be removed.\n", len(removed))
		return nil
	}

	cleaned := tr.Apply(secrets)
	w := env.NewWriter(envFile)
	if err := w.Write(cleaned); err != nil {
		return fmt.Errorf("writing env file: %w", err)
	}
	fmt.Printf("Removed %d placeholder key(s) from %s.\n", len(removed), envFile)
	return nil
}
