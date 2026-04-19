package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/env"
	"github.com/vaultline/vaultline/internal/envsign"
)

var signKey string
var signVerify string

func init() {
	signCmd := &cobra.Command{
		Use:   "sign",
		Short: "Sign or verify a .env file using HMAC-SHA256",
	}

	signCmd.PersistentFlags().StringVar(&signKey, "key", "", "HMAC signing key (required)")

	signCmd.AddCommand(&cobra.Command{
		Use:   "generate <env-file>",
		Short: "Generate a signature for a .env file",
		Args:  cobra.ExactArgs(1),
		RunE:  runSignGenerate,
	})

	verifyCmd := &cobra.Command{
		Use:   "verify <env-file>",
		Short: "Verify a .env file against a known signature",
		Args:  cobra.ExactArgs(1),
		RunE:  runSignVerify,
	}
	verifyCmd.Flags().StringVar(&signVerify, "sig", "", "Expected HMAC signature (required)")
	signCmd.AddCommand(verifyCmd)

	rootCmd.AddCommand(signCmd)
}

func runSignGenerate(cmd *cobra.Command, args []string) error {
	if signKey == "" {
		return fmt.Errorf("--key is required")
	}
	secrets, err := env.ReadFile(args[0])
	if err != nil {
		return fmt.Errorf("reading env file: %w", err)
	}
	s := envsign.New(signKey)
	sig, err := s.Sign(secrets)
	if err != nil {
		return fmt.Errorf("signing: %w", err)
	}
	fmt.Fprintln(os.Stdout, sig)
	return nil
}

func runSignVerify(cmd *cobra.Command, args []string) error {
	if signKey == "" {
		return fmt.Errorf("--key is required")
	}
	if signVerify == "" {
		return fmt.Errorf("--sig is required")
	}
	secrets, err := env.ReadFile(args[0])
	if err != nil {
		return fmt.Errorf("reading env file: %w", err)
	}
	s := envsign.New(signKey)
	if err := s.Verify(secrets, signVerify); err != nil {
		return fmt.Errorf("verification failed: %w", err)
	}
	fmt.Fprintln(os.Stdout, "signature valid")
	return nil
}
