package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/envcast"
	"github.com/vaultline/vaultline/internal/env"
)

var castFlags struct {
	envFile string
	types   []string
}

func init() {
	castCmd := &cobra.Command{
		Use:   "cast",
		Short: "Display secrets cast to typed values",
		RunE:  runCast,
	}
	castCmd.Flags().StringVar(&castFlags.envFile, "env-file", ".env", "Path to .env file")
	castCmd.Flags().StringArrayVar(&castFlags.types, "type", nil, "Key type pairs e.g. PORT=int (repeatable)")
	rootCmd.AddCommand(castCmd)
}

func runCast(cmd *cobra.Command, args []string) error {
	secrets, err := env.ReadFile(castFlags.envFile)
	if err != nil {
		return fmt.Errorf("cast: read env file: %w", err)
	}

	typeMap := map[string]string{}
	for _, pair := range castFlags.types {
		for i, ch := range pair {
			if ch == '=' {
				typeMap[pair[:i]] = pair[i+1:]
				break
			}
		}
	}

	c := envcast.New()
	out, errs := c.CastAll(secrets, typeMap)
	for _, e := range errs {
		fmt.Fprintf(os.Stderr, "warn: %v\n", e)
	}

	keys := make([]string, 0, len(out))
	for k := range out {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(cmd.OutOrStdout(), "%s = %v\n", k, out[k])
	}
	return nil
}
