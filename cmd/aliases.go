package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/envaliases"
	"github.com/vaultline/vaultline/internal/env"
)

var aliasesCmd = &cobra.Command{
	Use:   "aliases",
	Short: "Apply key aliases to a .env file and print the result",
	RunE:  runAliases,
}

func init() {
	aliasesCmd.Flags().StringP("input", "i", ".env", "source .env file")
	aliasesCmd.Flags().StringP("map", "m", "aliases.json", "JSON file containing alias map")
	aliasesCmd.Flags().Bool("strip", false, "strip aliased keys from output instead of adding them")
	RootCmd.AddCommand(aliasesCmd)
}

func runAliases(cmd *cobra.Command, _ []string) error {
	input, _ := cmd.Flags().GetString("input")
	mapFile, _ := cmd.Flags().GetString("map")
	strip, _ := cmd.Flags().GetBool("strip")

	secrets, err := env.ReadFile(input)
	if err != nil {
		return fmt.Errorf("aliases: read env: %w", err)
	}

	am, err := envaliases.LoadFile(mapFile)
	if err != nil {
		return fmt.Errorf("aliases: load map: %w", err)
	}

	a := envaliases.New(am)
	var result map[string]string
	if strip {
		result = a.Strip(secrets)
	} else {
		result = a.Apply(secrets)
	}

	keys := make([]string, 0, len(result))
	for k := range result {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	out := make(map[string]string, len(keys))
	for _, k := range keys {
		out[k] = result[k]
	}
	return enc.Encode(out)
}
