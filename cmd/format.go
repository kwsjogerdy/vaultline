package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/envformat"
)

var (
	formatSort   bool
	formatGroup  bool
	formatHeader string
	formatInput  string
)

func init() {
	formatCmd := &cobra.Command{
		Use:   "format",
		Short: "Render a JSON secret map as a formatted .env file",
		RunE:  runFormat,
	}
	formatCmd.Flags().StringVarP(&formatInput, "input", "i", "", "path to JSON secrets file (required)")
	formatCmd.Flags().BoolVar(&formatSort, "sort", true, "sort keys alphabetically")
	formatCmd.Flags().BoolVar(&formatGroup, "group", false, "insert blank lines between prefix groups")
	formatCmd.Flags().StringVar(&formatHeader, "header", "", "comment header written at the top of the output")
	_ = formatCmd.MarkFlagRequired("input")
	rootCmd.AddCommand(formatCmd)
}

func runFormat(cmd *cobra.Command, _ []string) error {
	data, err := os.ReadFile(formatInput)
	if err != nil {
		return fmt.Errorf("read input: %w", err)
	}
	var secrets map[string]string
	if err := json.Unmarshal(data, &secrets); err != nil {
		return fmt.Errorf("parse JSON: %w", err)
	}

	opts := envformat.Options{
		SortKeys:      formatSort,
		GroupByPrefix: formatGroup,
		Header:        formatHeader,
	}
	f := envformat.New(opts)
	if err := f.Write(cmd.OutOrStdout(), secrets); err != nil {
		return fmt.Errorf("format: %w", err)
	}
	return nil
}
