package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/decrypt"
)

var (
	encryptInput      string
	encryptPassphrase string
	encryptMode       string
)

func init() {
	encryptCmd := &cobra.Command{
		Use:   "encrypt",
		Short: "Encrypt or decrypt a secret value using AES-256-GCM",
		RunE:  runEncrypt,
	}
	encryptCmd.Flags().StringVarP(&encryptInput, "value", "v", "", "Value to encrypt or decrypt (required)")
	encryptCmd.Flags().StringVarP(&encryptPassphrase, "passphrase", "p", "", "Passphrase for key derivation (required)")
	encryptCmd.Flags().StringVarP(&encryptMode, "mode", "m", "encrypt", "Mode: encrypt or decrypt")
	_ = encryptCmd.MarkFlagRequired("value")
	_ = encryptCmd.MarkFlagRequired("passphrase")
	rootCmd.AddCommand(encryptCmd)
}

func runEncrypt(cmd *cobra.Command, args []string) error {
	key := decrypt.DeriveKey(encryptPassphrase)

	switch encryptMode {
	case "encrypt":
		result, err := decrypt.Encrypt(key, encryptInput)
		if err != nil {
			return fmt.Errorf("encryption failed: %w", err)
		}
		fmt.Fprintln(os.Stdout, result)
	case "decrypt":
		result, err := decrypt.Decrypt(key, encryptInput)
		if err != nil {
			return fmt.Errorf("decryption failed: %w", err)
		}
		fmt.Fprintln(os.Stdout, result)
	default:
		return fmt.Errorf("unknown mode %q: use 'encrypt' or 'decrypt'", encryptMode)
	}
	return nil
}
