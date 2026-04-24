package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/env"
	"github.com/vaultline/vaultline/internal/envpatch"
)

var (
	patchFile   string
	patchSet    []string
	patchDelete []string
	patchRename []string
)

func init() {
	patchCmd := &cobra.Command{
		Use:   "patch",
		Short: "Apply set/delete/rename operations to a .env file",
		RunE:  runPatch,
	}
	patchCmd.Flags().StringVarP(&patchFile, "file", "f", ".env", "target .env file")
	patchCmd.Flags().StringArrayVar(&patchSet, "set", nil, "KEY=VALUE pairs to set")
	patchCmd.Flags().StringArrayVar(&patchDelete, "delete", nil, "keys to delete")
	patchCmd.Flags().StringArrayVar(&patchRename, "rename", nil, "OLD=NEW key renames")
	rootCmd.AddCommand(patchCmd)
}

func runPatch(cmd *cobra.Command, _ []string) error {
	secrets, err := env.ReadFile(patchFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read %s: %w", patchFile, err)
	}
	if secrets == nil {
		secrets = map[string]string{}
	}

	var ops []envpatch.Op
	for _, kv := range patchSet {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("--set %q: expected KEY=VALUE", kv)
		}
		ops = append(ops, envpatch.Op{Type: envpatch.OpSet, Key: parts[0], Value: parts[1]})
	}
	for _, k := range patchDelete {
		ops = append(ops, envpatch.Op{Type: envpatch.OpDelete, Key: k})
	}
	for _, kv := range patchRename {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("--rename %q: expected OLD=NEW", kv)
		}
		ops = append(ops, envpatch.Op{Type: envpatch.OpRename, Key: parts[0], NewKey: parts[1]})
	}

	patcher, err := envpatch.New(ops)
	if err != nil {
		return err
	}

	out, res, err := patcher.Apply(secrets)
	if err != nil {
		return err
	}

	w := env.NewWriter(patchFile)
	if err := w.Write(out); err != nil {
		return fmt.Errorf("write %s: %w", patchFile, err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "patched %s: %d applied, %d skipped\n",
		patchFile, len(res.Applied), len(res.Skipped))
	return nil
}
