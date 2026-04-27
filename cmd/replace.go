package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/envreplace"
	"github.com/vaultline/vaultline/internal/env"
)

var replaceCmd = &cobra.Command{
	Use:   "replace",
	Short: "Find and replace substrings in secret values",
	RunE:  runReplace,
}

func init() {
	replaceCmd.Flags().String("file", ".env", "source .env file")
	replaceCmd.Flags().String("find", "", "substring to find (required)")
	replaceCmd.Flags().String("replace", "", "replacement string")
	replaceCmd.Flags().StringSlice("keys", nil, "limit replacement to specific keys")
	replaceCmd.Flags().String("out", "", "write result to file (default: stdout)")
	replaceCmd.Flags().String("format", "env", "output format: env or json")
	_ = replaceCmd.MarkFlagRequired("find")
	rootCmd.AddCommand(replaceCmd)
}

func runReplace(cmd *cobra.Command, _ []string) error {
	file, _ := cmd.Flags().GetString("file")
	find, _ := cmd.Flags().GetString("find")
	repl, _ := cmd.Flags().GetString("replace")
	keys, _ := cmd.Flags().GetStringSlice("keys")
	out, _ := cmd.Flags().GetString("out")
	format, _ := cmd.Flags().GetString("format")

	secrets, err := env.ReadFile(file)
	if err != nil {
		return fmt.Errorf("replace: read %s: %w", file, err)
	}

	r, err := envreplace.New([]envreplace.Rule{{Find: find, Replace: repl, Keys: keys}})
	if err != nil {
		return err
	}

	result := r.Apply(secrets)
	changed := r.Changed(secrets)

	var dest *os.File
	if out != "" {
		dest, err = os.Create(out)
		if err != nil {
			return fmt.Errorf("replace: create output: %w", err)
		}
		defer dest.Close()
	} else {
		dest = os.Stdout
	}

	switch format {
	case "json":
		enc := json.NewEncoder(dest)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	default:
		sortedKeys := make([]string, 0, len(result))
		for k := range result {
			sortedKeys = append(sortedKeys, k)
		}
		sort.Strings(sortedKeys)
		for _, k := range sortedKeys {
			fmt.Fprintf(dest, "%s=%s\n", k, result[k])
		}
	}

	if len(changed) > 0 {
		fmt.Fprintf(os.Stderr, "replaced in %d key(s): %v\n", len(changed), changed)
	}
	return nil
}
