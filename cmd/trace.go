package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/envtrace"
)

var (
	traceFile string
	traceKey  string
)

func init() {
	traceCmd := &cobra.Command{
		Use:   "trace",
		Short: "Show access history for secret keys",
		RunE:  runTrace,
	}
	traceCmd.Flags().StringVar(&traceFile, "file", ".vaultline-trace.json", "path to trace log file")
	traceCmd.Flags().StringVar(&traceKey, "key", "", "filter trace entries by key name")
	rootCmd.AddCommand(traceCmd)
}

func runTrace(cmd *cobra.Command, args []string) error {
	tr, err := envtrace.New(traceFile)
	if err != nil {
		return fmt.Errorf("trace: %w", err)
	}

	var entries []envtrace.Entry
	if traceKey != "" {
		entries = tr.ForKey(traceKey)
	} else {
		entries = tr.Entries()
	}

	if len(entries) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no trace entries found")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIME\tKEY\tOPERATION")
	for _, e := range entries {
		fmt.Fprintf(w, "%s\t%s\t%s\n", e.At.Format("2006-01-02 15:04:05"), e.Key, e.Operation)
	}
	return w.Flush()
}
