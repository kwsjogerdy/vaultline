package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/envpartition"
)

var (
	partitionPrefix  string
	partitionContain string
	partitionJSON    bool
)

func init() {
	partitionCmd := &cobra.Command{
		Use:   "partition",
		Short: "Split secrets into matched/unmatched groups based on a rule",
		RunE:  runPartition,
	}
	partitionCmd.Flags().StringVar(&partitionPrefix, "prefix", "", "Match keys with this prefix")
	partitionCmd.Flags().StringVar(&partitionContain, "contains", "", "Match keys containing this substring")
	partitionCmd.Flags().BoolVar(&partitionJSON, "json", false, "Output result as JSON")
	rootCmd.AddCommand(partitionCmd)
}

func runPartition(cmd *cobra.Command, args []string) error {
	if partitionPrefix == "" && partitionContain == "" {
		return fmt.Errorf("provide --prefix or --contains")
	}

	predicate := func(k, _ string) bool {
		if partitionPrefix != "" && !strings.HasPrefix(k, partitionPrefix) {
			return false
		}
		if partitionContain != "" && !strings.Contains(k, partitionContain) {
			return false
		}
		return true
	}

	p, err := envpartition.New(predicate)
	if err != nil {
		return err
	}

	secrets := map[string]string{}
	if err := json.NewDecoder(os.Stdin).Decode(&secrets); err != nil {
		return fmt.Errorf("reading secrets from stdin: %w", err)
	}

	res := p.Apply(secrets)

	if partitionJSON {
		out := map[string]map[string]string{
			"matched":   res.Matched,
			"unmatched": res.Unmatched,
		}
		return json.NewEncoder(os.Stdout).Encode(out)
	}

	fmt.Fprintf(os.Stdout, "matched (%d):\n", len(res.Matched))
	for _, k := range res.MatchedKeys() {
		fmt.Fprintf(os.Stdout, "  %s\n", k)
	}
	fmt.Fprintf(os.Stdout, "unmatched (%d):\n", len(res.Unmatched))
	for _, k := range res.UnmatchedKeys() {
		fmt.Fprintf(os.Stdout, "  %s\n", k)
	}
	return nil
}
