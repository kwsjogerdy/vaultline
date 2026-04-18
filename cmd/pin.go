package cmd

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/envpin"
)

var pinFile string

func init() {
	pinCmd := &cobra.Command{
		Use:   "pin",
		Short: "Manage secret version pins",
	}

	setCmd := &cobra.Command{
		Use:   "set <path> <version>",
		Short: "Pin a secret path to a specific version",
		Args:  cobra.ExactArgs(2),
		RunE:  runPinSet,
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all pinned secrets",
		RunE:  runPinList,
	}

	removeCmd := &cobra.Command{
		Use:   "remove <path>",
		Short: "Remove a pin for a secret path",
		Args:  cobra.ExactArgs(1),
		RunE:  runPinRemove,
	}

	pinCmd.PersistentFlags().StringVar(&pinFile, "pin-file", "vaultline.lock.json", "path to the pin lockfile")
	pinCmd.AddCommand(setCmd, listCmd, removeCmd)
	rootCmd.AddCommand(pinCmd)
}

func runPinSet(cmd *cobra.Command, args []string) error {
	version, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("version must be an integer: %w", err)
	}
	m := envpin.New(pinFile)
	lf, err := m.Load()
	if err != nil {
		return err
	}
	m.Set(lf, args[0], version)
	if err := m.Save(lf); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "pinned %s @ v%d\n", args[0], version)
	return nil
}

func runPinList(cmd *cobra.Command, _ []string) error {
	m := envpin.New(pinFile)
	lf, err := m.Load()
	if err != nil {
		return err
	}
	if len(lf.Pins) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no pins defined")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PATH\tVERSION\tPINNED AT")
	for _, p := range lf.Pins {
		fmt.Fprintf(w, "%s\t%d\t%s\n", p.Path, p.Version, p.PinnedAt.Format("2006-01-02 15:04:05"))
	}
	return w.Flush()
}

func runPinRemove(cmd *cobra.Command, args []string) error {
	m := envpin.New(pinFile)
	lf, err := m.Load()
	if err != nil {
		return err
	}
	m.Remove(lf, args[0])
	if err := m.Save(lf); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "removed pin for %s\n", args[0])
	return nil
}
