package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/envobfuscate"
)

var obfuscateCmd = &cobra.Command{
	Use:   "obfuscate",
	Short: "Obfuscate or deobfuscate secret values",
}

func init() {
	var (
		key     string
		input   string
		reverse bool
	)

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Obfuscate secrets from a JSON file (or deobfuscate with --reverse)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runObfuscate(key, input, reverse)
		},
	}
	runCmd.Flags().StringVar(&key, "key", "", "Obfuscation key (required)")
	runCmd.Flags().StringVar(&input, "input", "", "Path to JSON secrets file (required)")
	runCmd.Flags().BoolVar(&reverse, "reverse", false, "Deobfuscate instead of obfuscate")
	_ = runCmd.MarkFlagRequired("key")
	_ = runCmd.MarkFlagRequired("input")

	obfuscateCmd.AddCommand(runCmd)
	rootCmd.AddCommand(obfuscateCmd)
}

func runObfuscate(key, inputPath string, reverse bool) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read input: %w", err)
	}
	var secrets map[string]string
	if err := json.Unmarshal(data, &secrets); err != nil {
		return fmt.Errorf("parse input: %w", err)
	}
	o, err := envobfuscate.New(key)
	if err != nil {
		return err
	}
	var result map[string]string
	if reverse {
		result, err = o.Deobfuscate(secrets)
		if err != nil {
			return err
		}
	} else {
		result = o.Obfuscate(secrets)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(result)
}
