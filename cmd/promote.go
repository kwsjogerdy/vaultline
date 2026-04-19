package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"vaultline/internal/env"
	"vaultline/internal/envpromote"
)

var (
	promoteSrcFile  string
	promoteDstFile  string
	promotePrefix   string
	promoteOnlyKeys string
	promoteJSON     bool
)

func init() {
	promoteCmd := &cobra.Command{
		Use:   "promote",
		Short: "Promote secrets from one env file to another",
		RunE:  runPromote,
	}
	promoteCmd.Flags().StringVar(&promoteSrcFile, "src", "", "Source .env file (required)")
	promoteCmd.Flags().StringVar(&promoteDstFile, "dst", "", "Destination .env file (required)")
	promoteCmd.Flags().StringVar(&promotePrefix, "strip-prefix", "", "Key prefix to strip when promoting")
	promoteCmd.Flags().StringVar(&promoteOnlyKeys, "only", "", "Comma-separated list of keys to promote")
	promoteCmd.Flags().BoolVar(&promoteJSON, "json", false, "Output result as JSON")
	_ = promoteCmd.MarkFlagRequired("src")
	_ = promoteCmd.MarkFlagRequired("dst")
	rootCmd.AddCommand(promoteCmd)
}

func runPromote(cmd *cobra.Command, args []string) error {
	srcSecrets, err := env.ReadFile(promoteSrcFile)
	if err != nil {
		return fmt.Errorf("reading src: %w", err)
	}
	dstSecrets, _ := env.ReadFile(promoteDstFile) // dst may not exist yet
	if dstSecrets == nil {
		dstSecrets = map[string]string{}
	}

	var only []string
	if promoteOnlyKeys != "" {
		for _, k := range strings.Split(promoteOnlyKeys, ",") {
			only = append(only, strings.TrimSpace(k))
		}
	}

	p := envpromote.New(promoteSrcFile, promoteDstFile).WithPrefix(promotePrefix)
	merged, res := p.Promote(srcSecrets, dstSecrets, only)

	w := env.NewWriter(promoteDstFile)
	if err := w.Write(merged); err != nil {
		return fmt.Errorf("writing dst: %w", err)
	}

	if promoteJSON {
		return json.NewEncoder(os.Stdout).Encode(res)
	}
	fmt.Printf("Promoted %d key(s), skipped %d key(s) → %s\n", len(res.Copied), len(res.Skipped), promoteDstFile)
	return nil
}
