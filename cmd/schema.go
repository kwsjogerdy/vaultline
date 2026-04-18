package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/env"
	"github.com/vaultline/vaultline/internal/schema"
)

var (
	schemaFile  string
	schemaEnv   string
)

func init() {
	schemaCmd := &cobra.Command{
		Use:   "schema",
		Short: "Validate a .env file against a schema",
		RunE:  runSchema,
	}
	schemaCmd.Flags().StringVar(&schemaFile, "schema", "schema.json", "Path to schema JSON file")
	schemaCmd.Flags().StringVar(&schemaEnv, "env", ".env", "Path to .env file to validate")
	rootCmd.AddCommand(schemaCmd)
}

func runSchema(cmd *cobra.Command, args []string) error {
	s, err := schema.LoadFile(schemaFile)
	if err != nil {
		return fmt.Errorf("loading schema: %w", err)
	}

	secrets, err := env.ReadFile(schemaEnv)
	if err != nil {
		return fmt.Errorf("reading env file: %w", err)
	}

	if err := s.Validate(secrets); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Println("✓ schema validation passed")
	return nil
}
