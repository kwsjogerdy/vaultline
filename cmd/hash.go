package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/env"
	"github.com/vaultline/vaultline/internal/envhash"
)

var hashFile string
var hashJSON bool

func init() {
	hashCmd := &cobra.Command{
		Use:   "hash",
		Short: "Compute a stable fingerprint of a .env file",
		RunE:  runHash,
	}
	hashCmd.Flags().StringVarP(&hashFile, "file", "f", ".env", "Path to the .env file")
	hashCmd.Flags().BoolVar(&hashJSON, "json", false, "Output result as JSON")
	rootCmd.AddCommand(hashCmd)
}

func runHash(cmd *cobra.Command, args []string) error {
	secrets, err := env.ReadFile(hashFile)
	if err != nil {
		return fmt.Errorf("reading env file: %w", err)
	}

	h := envhash.New()
	result := h.Compute(secrets)

	if hashJSON {
		out := map[string]interface{}{
			"fingerprint": result.Fingerprint,
			"key_count":   result.KeyCount,
			"keys":        result.Keys,
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	}

	fmt.Printf("fingerprint : %s\n", result.Fingerprint)
	fmt.Printf("keys        : %d\n", result.KeyCount)
	return nil
}
