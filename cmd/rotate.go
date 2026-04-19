package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/env"
	"github.com/vaultline/vaultline/internal/envrotate"
)

var (
	rotateLedger  string
	rotateMaxAge  time.Duration
	rotateEnvFile string
)

func init() {
	rotateCmd := &cobra.Command{
		Use:   "rotate",
		Short: "Track and report secret rotation status",
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show keys due for rotation",
		RunE:  runRotateStatus,
	}
	statusCmd.Flags().StringVar(&rotateLedger, "ledger", ".rotate.json", "path to rotation ledger")
	statusCmd.Flags().StringVar(&rotateEnvFile, "env", ".env", "env file to inspect")
	statusCmd.Flags().DurationVar(&rotateMaxAge, "max-age", 30*24*time.Hour, "max age before rotation is due")

	markCmd := &cobra.Command{
		Use:   "mark [key...]",
		Short: "Mark keys as rotated now",
		Args:  cobra.MinimumNArgs(1),
		RunE:  runRotateMark,
	}
	markCmd.Flags().StringVar(&rotateLedger, "ledger", ".rotate.json", "path to rotation ledger")

	rotateCmd.AddCommand(statusCmd, markCmd)
	rootCmd.AddCommand(rotateCmd)
}

func runRotateStatus(cmd *cobra.Command, _ []string) error {
	secrets, err := env.ReadFile(rotateEnvFile)
	if err != nil {
		return fmt.Errorf("read env: %w", err)
	}
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	l, err := envrotate.New(rotateLedger)
	if err != nil {
		return err
	}
	due := l.DueForRotation(keys, rotateMaxAge)
	if len(due) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "All secrets are up to date.")
		return nil
	}
	fmt.Fprintf(cmd.OutOrStdout(), "%d key(s) due for rotation:\n", len(due))
	for _, k := range due {
		fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", k)
	}
	return nil
}

func runRotateMark(_ *cobra.Command, args []string) error {
	l, err := envrotate.New(rotateLedger)
	if err != nil {
		return err
	}
	for _, k := range args {
		r := l.Rotate(k)
		fmt.Fprintf(os.Stdout, "marked %s as rotated (version %d)\n", r.Key, r.Version)
	}
	return l.Save()
}
