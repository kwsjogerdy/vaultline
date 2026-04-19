package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/envfreeze"
)

var freezeFile string

func init() {
	freezeCmd := &cobra.Command{
		Use:   "freeze",
		Short: "Manage frozen secrets that are protected from sync overwrites",
	}

	freezePinCmd := &cobra.Command{
		Use:   "add <key>",
		Short: "Freeze a secret key",
		Args:  cobra.ExactArgs(1),
		RunE:  runFreezeAdd,
	}

	freezeRemoveCmd := &cobra.Command{
		Use:   "remove <key>",
		Short: "Unfreeze a secret key",
		Args:  cobra.ExactArgs(1),
		RunE:  runFreezeRemove,
	}

	freezeListCmd := &cobra.Command{
		Use:   "list",
		Short: "List all frozen keys",
		RunE:  runFreezeList,
	}

	for _, sub := range []*cobra.Command{freezePinCmd, freezeRemoveCmd, freezeListCmd} {
		sub.Flags().StringVar(&freezeFile, "freeze-file", ".vaultline-freeze.json", "Path to freeze ledger file")
		freezeCmd.AddCommand(sub)
	}

	rootCmd.AddCommand(freezeCmd)
}

func runFreezeAdd(cmd *cobra.Command, args []string) error {
	f, err := envfreeze.New(freezeFile)
	if err != nil {
		return err
	}
	if err := f.Freeze(args[0]); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "frozen: %s\n", args[0])
	return nil
}

func runFreezeRemove(cmd *cobra.Command, args []string) error {
	f, err := envfreeze.New(freezeFile)
	if err != nil {
		return err
	}
	if err := f.Unfreeze(args[0]); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "unfrozen: %s\n", args[0])
	return nil
}

func runFreezeList(cmd *cobra.Command, args []string) error {
	f, err := envfreeze.New(freezeFile)
	if err != nil {
		return err
	}
	entries := f.List()
	if len(entries) == 0 {
		fmt.Fprintln(os.Stdout, "no frozen keys")
		return nil
	}
	for _, e := range entries {
		fmt.Fprintf(os.Stdout, "%-30s  frozen at %s\n", e.Key, e.FrozenAt.Format("2006-01-02 15:04:05"))
	}
	return nil
}
