package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vaultline/vaultline/internal/config"
	"github.com/vaultline/vaultline/internal/template"
	"github.com/vaultline/vaultline/internal/vault"
)

var (
	tmplFile   string
	tmplOutput string
)

func init() {
	tmplCmd := &cobra.Command{
		Use:   "template",
		Short: "Render a .env file from a template using Vault secrets",
		RunE:  runTemplate,
	}
	tmplCmd.Flags().StringVarP(&tmplFile, "template", "t", ".env.tmpl", "Path to the template file")
	tmplCmd.Flags().StringVarP(&tmplOutput, "output", "o", ".env", "Path to write the rendered output")
	RootCmd.AddCommand(tmplCmd)
}

func runTemplate(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	secrets, err := client.GetSecrets(cmd.Context(), cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("vault: %w", err)
	}

	r := template.New(tmplFile)
	output, err := r.Render(secrets)
	if err != nil {
		return fmt.Errorf("render: %w", err)
	}

	if err := os.WriteFile(tmplOutput, []byte(output), 0600); err != nil {
		return fmt.Errorf("write %q: %w", tmplOutput, err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "rendered template → %s\n", tmplOutput)
	return nil
}
