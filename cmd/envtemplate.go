package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/env"
	"github.com/vaultline/vaultline/internal/envtemplate"
)

var (
	envtmplInput  string
	envtmplOutput string
	envtmplSource string
)

func init() {
	envtmplCmd := &cobra.Command{
		Use:   "envtemplate",
		Short: "Render a template file using secrets from a .env source",
		RunE:  runEnvTemplate,
	}
	envtmplCmd.Flags().StringVar(&envtmplInput, "template", "", "path to template file (required)")
	envtmplCmd.Flags().StringVar(&envtmplOutput, "output", "", "path to write rendered output (default: stdout)")
	envtmplCmd.Flags().StringVar(&envtmplSource, "source", ".env", "path to .env file supplying secrets")
	_ = envtmplCmd.MarkFlagRequired("template")
	rootCmd.AddCommand(envtmplCmd)
}

func runEnvTemplate(cmd *cobra.Command, _ []string) error {
	secrets, err := env.ReadFile(envtmplSource)
	if err != nil {
		return fmt.Errorf("read source: %w", err)
	}

	tmplBytes, err := os.ReadFile(envtmplInput)
	if err != nil {
		return fmt.Errorf("read template: %w", err)
	}

	r := envtemplate.New(secrets)
	result, err := r.Render(string(tmplBytes))
	if err != nil {
		return fmt.Errorf("render: %w", err)
	}

	if envtmplOutput == "" {
		fmt.Fprint(cmd.OutOrStdout(), result)
		return nil
	}

	if err := os.WriteFile(envtmplOutput, []byte(result), 0o600); err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "rendered output written to %s\n", envtmplOutput)
	return nil
}
