package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/config"
	"github.com/vaultline/vaultline/internal/envclone"
	"github.com/vaultline/vaultline/internal/vault"
)

var (
	cloneOnly  []string
	cloneRemap []string
)

func init() {
	cloneCmd := &cobra.Command{
		Use:   "clone <src-path> <dst-path>",
		Short: "Clone secrets from one Vault path to another",
		Args:  cobra.ExactArgs(2),
		RunE:  runClone,
	}
	cloneCmd.Flags().StringSliceVar(&cloneOnly, "only", nil, "comma-separated keys to copy")
	cloneCmd.Flags().StringArrayVar(&cloneRemap, "remap", nil, "key remaps as OLD=NEW (repeatable)")
	rootCmd.AddCommand(cloneCmd)
}

func runClone(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	remap := make(map[string]string)
	for _, r := range cloneRemap {
		parts := strings.SplitN(r, "=", 2)
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "warning: ignoring invalid remap %q\n", r)
			continue
		}
		remap[parts[0]] = parts[1]
	}

	cloner := envclone.New(client, client)
	n, err := cloner.Clone(args[0], args[1], envclone.Options{
		Only:  cloneOnly,
		Remap: remap,
	})
	if err != nil {
		return err
	}

	fmt.Printf("cloned %d secret(s) from %q to %q\n", n, args[0], args[1])
	return nil
}
