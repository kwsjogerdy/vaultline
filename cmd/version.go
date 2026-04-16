package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/vaultline/vaultline/internal/version"
)

var jsonOutput bool

func init() {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print vaultline version information",
		RunE:  runVersion,
	}
	versionCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output version info as JSON")
	rootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) error {
	info := version.Get()
	if jsonOutput {
		fmt.Fprintf(cmd.OutOrStdout(),
			`{"version":%q,"commit":%q,"buildDate":%q,"goVersion":%q}`+"\n",
			info.Version, info.Commit, info.BuildDate, info.GoVersion)
		return nil
	}
	fmt.Fprintln(cmd.OutOrStdout(), info.String())
	return nil
}
