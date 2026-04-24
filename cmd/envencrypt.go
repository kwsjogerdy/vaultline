package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/env"
	"github.com/vaultline/vaultline/internal/envencrypt"
)

var (
	encryptFile       string
	encryptPassphrase string
	encryptKeys       []string
	encryptDecrypt    bool
)

func init() {
	encryptCmd := &cobra.Command{
		Use:   "envencrypt",
		Short: "Encrypt or decrypt values inside a .env file",
		RunE:  runEnvEncrypt,
	}
	encryptCmd.Flags().StringVarP(&encryptFile, "file", "f", ".env", "path to the .env file")
	encryptCmd.Flags().StringVarP(&encryptPassphrase, "passphrase", "p", "", "encryption passphrase (required)")
	encryptCmd.Flags().StringSliceVarP(&encryptKeys, "keys", "k", nil, "keys to encrypt/decrypt (default: all)")
	encryptCmd.Flags().BoolVarP(&encryptDecrypt, "decrypt", "d", false, "decrypt instead of encrypt")
	_ = encryptCmd.MarkFlagRequired("passphrase")
	rootCmd.AddCommand(encryptCmd)
}

func runEnvEncrypt(cmd *cobra.Command, _ []string) error {
	secrets, err := env.ReadFile(encryptFile)
	if err != nil {
		return fmt.Errorf("envencrypt: read file: %w", err)
	}

	e, err := envencrypt.New(encryptPassphrase, encryptKeys)
	if err != nil {
		return err
	}

	var result map[string]string
	if encryptDecrypt {
		result, err = e.Decrypt(secrets)
	} else {
		result, err = e.Encrypt(secrets)
	}
	if err != nil {
		return err
	}

	w := env.NewWriter(encryptFile)
	if err := w.Write(result); err != nil {
		return fmt.Errorf("envencrypt: write file: %w", err)
	}

	action := "Encrypted"
	if encryptDecrypt {
		action = "Decrypted"
	}
	target := "all keys"
	if len(encryptKeys) > 0 {
		target = strings.Join(encryptKeys, ", ")
	}
	fmt.Fprintf(os.Stdout, "%s %s in %s\n", action, target, encryptFile)
	return nil
}
