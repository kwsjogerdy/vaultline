package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/env"
	"github.com/vaultline/vaultline/internal/envprefix"
)

var (
	prefixAdd    string
	prefixRemove string
	prefixOld    string
	prefixNew    string
	prefixInput  string
	prefixOutput string
)

func init() {
	prefixCmd := &cobra.Command{
		Use:   "prefix",
		Short: "Add, remove, or replace key prefixes in a .env file",
		RunE:  runPrefix,
	}
	prefixCmd.Flags().StringVar(&prefixAdd, "add", "", "prefix to add to all keys")
	prefixCmd.Flags().StringVar(&prefixRemove, "remove", "", "prefix to strip from all keys")
	prefixCmd.Flags().StringVar(&prefixOld, "replace-old", "", "prefix to replace (use with --replace-new)")
	prefixCmd.Flags().StringVar(&prefixNew, "replace-new", "", "replacement prefix (use with --replace-old)")
	prefixCmd.Flags().StringVar(&prefixInput, "input", ".env", "input .env file")
	prefixCmd.Flags().StringVar(&prefixOutput, "output", "", "output file (default: stdout)")
	rootCmd.AddCommand(prefixCmd)
}

func runPrefix(cmd *cobra.Command, args []string) error {
	secrets, err := env.ReadFile(prefixInput)
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	var result map[string]string

	if prefixOld != "" || prefixNew != "" {
		result = envprefix.ReplacePrefix(secrets, prefixOld, prefixNew)
	} else {
		tr, err := envprefix.New(prefixAdd, prefixRemove)
		if err != nil {
			return err
		}
		result = tr.Apply(secrets)
	}

	w := env.NewWriter(prefixOutput)
	if prefixOutput == "" {
		w = env.NewWriter("")
		for k, v := range result {
			fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
		}
		return nil
	}
	_ = w
	return env.NewWriter(prefixOutput).Write(result)
}
