package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/env"
	"github.com/vaultline/vaultline/internal/envfence"
)

func init() {
	fenceCmd := &cobra.Command{
		Use:   "fence",
		Short: "Filter secrets using an allow or deny key policy",
		RunE:  runFence,
	}
	fenceCmd.Flags().String("file", ".env", "Source .env file")
	fenceCmd.Flags().String("policy", "allow", "Policy: allow or deny")
	fenceCmd.Flags().StringSlice("keys", nil, "Comma-separated list of keys")
	fenceCmd.Flags().Bool("show-blocked", false, "Print blocked keys instead of allowed")
	rootCmd.AddCommand(fenceCmd)
}

func runFence(cmd *cobra.Command, _ []string) error {
	file, _ := cmd.Flags().GetString("file")
	policyStr, _ := cmd.Flags().GetString("policy")
	keys, _ := cmd.Flags().GetStringSlice("keys")
	showBlocked, _ := cmd.Flags().GetBool("show-blocked")

	secrets, err := env.ReadFile(file)
	if err != nil {
		return fmt.Errorf("fence: read %s: %w", file, err)
	}

	f, err := envfence.New(envfence.Policy(policyStr), keys)
	if err != nil {
		return err
	}

	if showBlocked {
		blocked := f.Blocked(secrets)
		if len(blocked) == 0 {
			fmt.Fprintln(os.Stdout, "no blocked keys")
			return nil
		}
		sort.Strings(blocked)
		for _, k := range blocked {
			fmt.Fprintln(os.Stdout, k)
		}
		return nil
	}

	out := f.Apply(secrets)
	sortedKeys := make([]string, 0, len(out))
	for k := range out {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	for _, k := range sortedKeys {
		fmt.Fprintf(os.Stdout, "%s=%s\n", k, strings.TrimSpace(out[k]))
	}
	return nil
}
