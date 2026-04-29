package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/env"
	"github.com/vaultline/vaultline/internal/envduration"
)

func init() {
	var (
		envFile string
		keys    []string
	)

	cmd := &cobra.Command{
		Use:   "duration",
		Short: "Parse and display duration-valued secrets from a .env file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDuration(envFile, keys)
		},
	}

	cmd.Flags().StringVar(&envFile, "file", ".env", "path to .env file")
	cmd.Flags().StringSliceVar(&keys, "keys", nil, "specific keys to parse (error on failure); omit to parse all")

	rootCmd.AddCommand(cmd)
}

func runDuration(envFile string, keys []string) error {
	secrets, err := env.ReadFile(envFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", envFile)
		}
		return fmt.Errorf("reading %s: %w", envFile, err)
	}

	p := envduration.New(keys...)
	results, err := p.Parse(secrets)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Println("no duration values found")
		return nil
	}

	width := 0
	for _, r := range results {
		if len(r.Key) > width {
			width = len(r.Key)
		}
	}

	for _, r := range results {
		padding := strings.Repeat(" ", width-len(r.Key))
		fmt.Printf("  %s%s  =  %v  (raw: %s)\n", r.Key, padding, r.Duration, r.Raw)
	}
	return nil
}
