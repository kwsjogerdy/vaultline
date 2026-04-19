package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/env"
	"github.com/vaultline/vaultline/internal/envcompare"
)

var compareCmd = &cobra.Command{
	Use:   "compare <left-env> <right-env>",
	Short: "Compare secrets between two .env files",
	Args:  cobra.ExactArgs(2),
	RunE:  runCompare,
}

var compareMaskValues bool

func init() {
	compareCmd.Flags().BoolVar(&compareMaskValues, "mask", true, "Mask secret values in output")
	rootCmd.AddCommand(compareCmd)
}

func runCompare(cmd *cobra.Command, args []string) error {
	leftPath, rightPath := args[0], args[1]

	left, err := env.ReadFile(leftPath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", leftPath, err)
	}
	right, err := env.ReadFile(rightPath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", rightPath, err)
	}

	c := envcompare.New(leftPath, rightPath)
	report := c.Compare(left, right)

	for _, res := range report.Results {
		lv, rv := res.Left, res.Right
		if compareMaskValues {
			if lv != "" {
				lv = "***"
			}
			if rv != "" {
				rv = "***"
			}
		}
		switch res.Status {
		case "match":
			fmt.Fprintf(os.Stdout, "  = %s\n", res.Key)
		case "differ":
			fmt.Fprintf(os.Stdout, "  ~ %s  (%s | %s)\n", res.Key, lv, rv)
		case "left_only":
			fmt.Fprintf(os.Stdout, "  < %s  (only in %s)\n", res.Key, report.Left)
		case "right_only":
			fmt.Fprintf(os.Stdout, "  > %s  (only in %s)\n", res.Key, report.Right)
		}
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, report.Summary())
	return nil
}
