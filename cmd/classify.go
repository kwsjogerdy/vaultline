package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/envclassify"
	"github.com/vaultline/vaultline/internal/env"
)

var classifyFormat string

func init() {
	classifyCmd := &cobra.Command{
		Use:   "classify",
		Short: "Classify secret keys into semantic categories",
		RunE:  runClassify,
	}
	classifyCmd.Flags().StringVarP(&classifyFormat, "format", "f", "text", "Output format: text or json")
	classifyCmd.Flags().StringP("file", "e", ".env", "Path to .env file")
	rootCmd.AddCommand(classifyCmd)
}

func runClassify(cmd *cobra.Command, _ []string) error {
	filePath, _ := cmd.Flags().GetString("file")

	secrets, err := env.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("reading env file: %w", err)
	}

	c := envclassify.New()
	results := c.Classify(secrets)

	switch classifyFormat {
	case "json":
		return json.NewEncoder(os.Stdout).Encode(results)
	default:
		grouped := envclassify.ByCategory(results)
		categories := []envclassify.Category{
			envclassify.CategoryCredential,
			envclassify.CategoryEndpoint,
			envclassify.CategoryFeatureFlag,
			envclassify.CategoryIdentifier,
			envclassify.CategoryUnknown,
		}
		for _, cat := range categories {
			items := grouped[cat]
			if len(items) == 0 {
				continue
			}
			fmt.Fprintf(os.Stdout, "[%s]\n", cat)
			for _, r := range items {
				fmt.Fprintf(os.Stdout, "  %s\n", r.Key)
			}
		}
	}
	return nil
}
