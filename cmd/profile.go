package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/profile"
)

var profileDir string

func init() {
	profileCmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage named sync profiles",
	}

	saveCmd := &cobra.Command{
		Use:   "save <name>",
		Short: "Save a profile",
		Args:  cobra.ExactArgs(1),
		RunE:  runProfileSave,
	}
	saveCmd.Flags().String("vault-path", "", "Vault secret path")
	saveCmd.Flags().String("output", ".env", "Output .env file")
	saveCmd.Flags().StringSlice("keys", nil, "Keys to include")
	saveCmd.Flags().StringSlice("exclude", nil, "Keys to exclude")

	listCmd := &cobra.Command{Use: "list", Short: "List profiles", RunE: runProfileList}
	deleteCmd := &cobra.Command{Use: "delete <name>", Short: "Delete a profile", Args: cobra.ExactArgs(1), RunE: runProfileDelete}

	profileCmd.PersistentFlags().StringVar(&profileDir, "profile-dir", defaultProfileDir(), "Profile storage directory")
	profileCmd.AddCommand(saveCmd, listCmd, deleteCmd)
	rootCmd.AddCommand(profileCmd)
}

func defaultProfileDir() string {
	home, _ := os.UserHomeDir()
	return home + "/.vaultline/profiles"
}

func runProfileSave(cmd *cobra.Command, args []string) error {
	vaultPath, _ := cmd.Flags().GetString("vault-path")
	output, _ := cmd.Flags().GetString("output")
	keys, _ := cmd.Flags().GetStringSlice("keys")
	exclude, _ := cmd.Flags().GetStringSlice("exclude")
	p := profile.Profile{Name: args[0], VaultPath: vaultPath, OutputFile: output, Keys: keys, Exclude: exclude}
	if err := profile.New(profileDir).Save(p); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "profile %q saved\n", args[0])
	return nil
}

func runProfileList(cmd *cobra.Command, _ []string) error {
	names, err := profile.New(profileDir).List()
	if err != nil {
		return err
	}
	if len(names) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no profiles found")
		return nil
	}
	fmt.Fprintln(cmd.OutOrStdout(), strings.Join(names, "\n"))
	return nil
}

func runProfileDelete(cmd *cobra.Command, args []string) error {
	if err := profile.New(profileDir).Delete(args[0]); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "profile %q deleted\n", args[0])
	return nil
}
