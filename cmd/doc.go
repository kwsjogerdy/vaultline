package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/envdoc"
)

var (
	docFormat  string
	docFile    string
)

func init() {
	docCmd := &cobra.Command{
		Use:   "doc",
		Short: "Generate documentation for secrets from an annotation file",
		RunE:  runDoc,
	}
	docCmd.Flags().StringVarP(&docFormat, "format", "f", "text", "Output format: text or markdown")
	docCmd.Flags().StringVarP(&docFile, "file", "i", "envdoc.json", "Path to the annotation JSON file")
	rootCmd.AddCommand(docCmd)
}

func runDoc(cmd *cobra.Command, args []string) error {
	d, err := envdoc.LoadFile(docFile)
	if err != nil {
		return fmt.Errorf("doc: %w", err)
	}
	if len(d.Keys()) == 0 {
		fmt.Fprintln(os.Stderr, "doc: no entries found in annotation file")
		return nil
	}
	if err := d.Write(os.Stdout, docFormat); err != nil {
		return fmt.Errorf("doc: %w", err)
	}
	return nil
}
