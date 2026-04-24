package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/envlabel"
)

var (
	labelFile   string
	labelFilter string
)

func init() {
	labelCmd := &cobra.Command{
		Use:   "label",
		Short: "Manage metadata labels for secret keys",
	}

	setCmd := &cobra.Command{
		Use:   "set <secret-key> <label-key>=<label-value>",
		Short: "Attach a label to a secret key",
		Args:  cobra.ExactArgs(2),
		RunE:  runLabelSet,
	}
	setCmd.Flags().StringVar(&labelFile, "file", ".vaultline-labels.json", "label store file")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List secret keys matching a label filter (key=value)",
		Args:  cobra.NoArgs,
		RunE:  runLabelList,
	}
	listCmd.Flags().StringVar(&labelFile, "file", ".vaultline-labels.json", "label store file")
	listCmd.Flags().StringVar(&labelFilter, "match", "", "filter as key=value (required)")
	_ = listCmd.MarkFlagRequired("match")

	labelCmd.AddCommand(setCmd, listCmd)
	rootCmd.AddCommand(labelCmd)
}

func runLabelSet(cmd *cobra.Command, args []string) error {
	secretKey := args[0]
	parts := strings.SplitN(args[1], "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("label must be in key=value format")
	}
	labelKey, labelValue := parts[0], parts[1]

	l, err := envlabel.New(labelFile)
	if err != nil {
		return err
	}
	l.Set(secretKey, labelKey, labelValue)
	if err := l.Save(); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "labeled %s [%s=%s]\n", secretKey, labelKey, labelValue)
	return nil
}

func runLabelList(cmd *cobra.Command, _ []string) error {
	parts := strings.SplitN(labelFilter, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("--match must be in key=value format")
	}
	labelKey, labelValue := parts[0], parts[1]

	l, err := envlabel.New(labelFile)
	if err != nil {
		return err
	}
	keys := l.Filter(labelKey, labelValue)
	if len(keys) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no secrets match the given label")
		return nil
	}
	for _, k := range keys {
		fmt.Fprintln(os.Stdout, k)
	}
	return nil
}
