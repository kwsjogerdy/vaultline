package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/envreorder"
)

var reorderCmd = &cobra.Command{
	Use:   "reorder",
	Short: "Reorder secrets according to a defined key order",
	RunE:  runReorder,
}

func init() {
	reorderCmd.Flags().StringSlice("keys", nil, "Comma-separated ordered key list (required keys appear first)")
	reorderCmd.Flags().Bool("fallback", true, "Sort remaining keys alphabetically")
	reorderCmd.Flags().String("input", "", "JSON file of secrets to reorder (default: stdin)")
	RootCmd.AddCommand(reorderCmd)
}

func runReorder(cmd *cobra.Command, _ []string) error {
	keys, _ := cmd.Flags().GetStringSlice("keys")
	fallback, _ := cmd.Flags().GetBool("fallback")
	input, _ := cmd.Flags().GetString("input")

	r, err := envreorder.New(keys, fallback)
	if err != nil {
		return fmt.Errorf("reorder: %w", err)
	}

	var raw []byte
	if input != "" {
		raw, err = os.ReadFile(input)
		if err != nil {
			return fmt.Errorf("reorder: reading input: %w", err)
		}
	} else {
		raw, err = io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("reorder: reading stdin: %w", err)
		}
	}

	var secrets map[string]string
	if err := json.Unmarshal(raw, &secrets); err != nil {
		return fmt.Errorf("reorder: parsing JSON: %w", err)
	}

	kvs := r.Apply(secrets)

	w := cmd.OutOrStdout()
	for _, kv := range kvs {
		fmt.Fprintf(w, "%s=%s\n", kv.Key, strings.ReplaceAll(kv.Value, "\n", "\\n"))
	}
	return nil
}
