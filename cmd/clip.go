package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"vaultline/internal/env"
	"vaultline/internal/envclip"
)

var (
	clipEnvFile    string
	clipClearAfter time.Duration
)

func init() {
	clipCmd := &cobra.Command{
		Use:   "clip <KEY>",
		Short: "Copy a secret value to the clipboard",
		Args:  cobra.ExactArgs(1),
		RunE:  runClip,
	}
	clipCmd.Flags().StringVarP(&clipEnvFile, "file", "f", ".env", "Path to .env file")
	clipCmd.Flags().DurationVar(&clipClearAfter, "clear-after", 30*time.Second, "Auto-clear clipboard after duration (0 to disable)")
	rootCmd.AddCommand(clipCmd)
}

func runClip(cmd *cobra.Command, args []string) error {
	key := args[0]

	secrets, err := env.ReadFile(clipEnvFile)
	if err != nil {
		return fmt.Errorf("read env file: %w", err)
	}

	c := envclip.New(clipClearAfter)
	if err := c.Copy(secrets, key); err != nil {
		return err
	}

	msg := fmt.Sprintf("✓ Copied %s to clipboard", key)
	if clipClearAfter > 0 {
		msg += fmt.Sprintf(" (clears in %s)", clipClearAfter)
	}
	fmt.Fprintln(os.Stdout, msg)
	return nil
}
