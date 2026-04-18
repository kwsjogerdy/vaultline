package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/env"
	"github.com/vaultline/vaultline/internal/envsnapshot"
)

var snapshotDir string

func init() {
	snapshotCmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Manage point-in-time snapshots of your .env file",
	}

	saveCmd := &cobra.Command{
		Use:   "save [label]",
		Short: "Save current .env as a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE:  runSnapshotSave,
	}
	saveCmd.Flags().StringVar(&snapshotDir, "snapshot-dir", ".vaultline/snapshots", "directory to store snapshots")
	saveCmd.Flags().StringP("env-file", "e", ".env", "source .env file")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List saved snapshots",
		RunE:  runSnapshotList,
	}
	listCmd.Flags().StringVar(&snapshotDir, "snapshot-dir", ".vaultline/snapshots", "directory to store snapshots")

	snapshotCmd.AddCommand(saveCmd, listCmd)
	rootCmd.AddCommand(snapshotCmd)
}

func runSnapshotSave(cmd *cobra.Command, args []string) error {
	label := args[0]
	envFile, _ := cmd.Flags().GetString("env-file")
	dir, _ := cmd.Flags().GetString("snapshot-dir")

	secrets, err := env.ReadFile(envFile)
	if err != nil {
		return fmt.Errorf("reading env file: %w", err)
	}

	m := envsnapshot.New(dir)
	if err := m.Save(label, secrets); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "Snapshot %q saved (%d keys)\n", label, len(secrets))
	return nil
}

func runSnapshotList(cmd *cobra.Command, args []string) error {
	dir, _ := cmd.Flags().GetString("snapshot-dir")
	m := envsnapshot.New(dir)
	list, err := m.List()
	if err != nil {
		return err
	}
	if len(list) == 0 {
		fmt.Println("No snapshots found.")
		return nil
	}
	for _, s := range list {
		fmt.Printf("%-20s  %s  (%d keys)\n", s.Label, s.CreatedAt.Format("2006-01-02 15:04:05"), len(s.Secrets))
	}
	return nil
}
