package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/envsanitize"
)

var sanitizeCmd = &cobra.Command{
	Use:   "sanitize",
	Short: "Sanitize secret values from a JSON file",
	RunE:  runSanitize,
}

var (
	sanitizeInput      string
	sanitizeNoTrim     bool
	sanitizeKeepNL     bool
	sanitizeKeepNulls  bool
)

func init() {
	rootCmd.AddCommand(sanitizeCmd)
	sanitizeCmd.Flags().StringVarP(&sanitizeInput, "input", "i", "", "path to JSON secrets file (required)")
	sanitizeCmd.Flags().BoolVar(&sanitizeNoTrim, "no-trim", false, "disable whitespace trimming")
	sanitizeCmd.Flags().BoolVar(&sanitizeKeepNL, "keep-newlines", false, "preserve newline characters in values")
	sanitizeCmd.Flags().BoolVar(&sanitizeKeepNulls, "keep-nulls", false, "preserve null bytes in values")
	_ = sanitizeCmd.MarkFlagRequired("input")
}

func runSanitize(cmd *cobra.Command, _ []string) error {
	data, err := os.ReadFile(sanitizeInput)
	if err != nil {
		return fmt.Errorf("sanitize: read input: %w", err)
	}

	var secrets map[string]string
	if err := json.Unmarshal(data, &secrets); err != nil {
		return fmt.Errorf("sanitize: parse JSON: %w", err)
	}

	opts := envsanitize.Options{
		TrimSpace:     !sanitizeNoTrim,
		StripNewlines: !sanitizeKeepNL,
		StripNulls:    !sanitizeKeepNulls,
	}

	s, err := envsanitize.New(opts)
	if err != nil {
		return fmt.Errorf("sanitize: build sanitizer: %w", err)
	}

	out := s.Apply(secrets)

	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		return fmt.Errorf("sanitize: encode output: %w", err)
	}
	return nil
}
