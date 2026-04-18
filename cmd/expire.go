package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/envexpire"
)

var (
	expireFile string
	expireTTL  time.Duration
)

func init() {
	expireCmd := &cobra.Command{
		Use:   "expire",
		Short: "Manage secret expiry tracking",
	}

	setCmd := &cobra.Command{
		Use:   "set <key>",
		Short: "Set expiry TTL for a secret key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runExpireSet(args[0])
		},
	}
	setCmd.Flags().StringVar(&expireFile, "file", ".vaultline-expire.json", "expiry store file")
	setCmd.Flags().DurationVar(&expireTTL, "ttl", 24*time.Hour, "time-to-live for the secret")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List expired secret keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runExpireList()
		},
	}
	listCmd.Flags().StringVar(&expireFile, "file", ".vaultline-expire.json", "expiry store file")

	expireCmd.AddCommand(setCmd, listCmd)
	rootCmd.AddCommand(expireCmd)
}

func runExpireSet(key string) error {
	s, err := envexpire.New(expireFile)
	if err != nil {
		return err
	}
	s.Set(key, expireTTL)
	if err := s.Save(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "set expiry for %q: %s\n", key, expireTTL)
	return nil
}

func runExpireList() error {
	s, err := envexpire.New(expireFile)
	if err != nil {
		return err
	}
	expired := s.Expired()
	if len(expired) == 0 {
		fmt.Println("no expired secrets")
		return nil
	}
	fmt.Println("expired secrets:")
	for _, k := range expired {
		fmt.Printf("  - %s\n", k)
	}
	return nil
}
