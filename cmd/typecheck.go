package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/env"
	"github.com/vaultline/vaultline/internal/envtypecheck"
)

var (
	typecheckFile   string
	typecheckStrict bool
	typecheckEnforce []string
)

func init() {
	typecheckCmd := &cobra.Command{
		Use:   "typecheck",
		Short: "Infer and validate value types in a .env file",
		RunE:  runTypecheck,
	}
	typecheckCmd.Flags().StringVarP(&typecheckFile, "file", "f", ".env", "Path to .env file")
	typecheckCmd.Flags().BoolVar(&typecheckStrict, "strict", false, "Exit non-zero if any type mismatch found")
	typecheckCmd.Flags().StringArrayVar(&typecheckEnforce, "enforce", nil, "Enforce a type: KEY=type (e.g. PORT=int)")
	rootCmd.AddCommand(typecheckCmd)
}

func runTypecheck(cmd *cobra.Command, args []string) error {
	secrets, err := env.ReadFile(typecheckFile)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	enforce := map[string]envtypecheck.Type{}
	for _, e := range typecheckEnforce {
		for i := 0; i < len(e); i++ {
			if e[i] == '=' {
				enforce[e[:i]] = envtypecheck.Type(e[i+1:])
				break
			}
		}
	}

	checker := envtypecheck.New(enforce)
	results := checker.Check(secrets)

	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})

	hasErr := false
	for _, r := range results {
		if r.Error != nil {
			fmt.Fprintf(os.Stderr, "MISMATCH  %-30s %-10s %s\n", r.Key, r.Type, r.Error)
			hasErr = true
		} else {
			fmt.Printf("OK        %-30s %s\n", r.Key, r.Type)
		}
	}

	if typecheckStrict && hasErr {
		return fmt.Errorf("type check failed: mismatches detected")
	}
	return nil
}
