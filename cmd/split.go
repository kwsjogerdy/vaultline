package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/envsplit"
)

var splitCmd = &cobra.Command{
	Use:   "split",
	Short: "Split secrets into named buckets by prefix rules",
	RunE:  runSplit,
}

func init() {
	splitCmd.Flags().StringArray("rule", nil, "Routing rule in bucket=PREFIX format (repeatable)")
	splitCmd.Flags().String("input", "", "JSON file of secrets to split (default: stdin)")
	splitCmd.Flags().Bool("show-unrouted", true, "Print unrouted keys")
	rootCmd.AddCommand(splitCmd)
}

func runSplit(cmd *cobra.Command, _ []string) error {
	ruleStrs, _ := cmd.Flags().GetStringArray("rule")
	inputFile, _ := cmd.Flags().GetString("input")
	showUnrouted, _ := cmd.Flags().GetBool("show-unrouted")

	if len(ruleStrs) == 0 {
		return fmt.Errorf("at least one --rule bucket=PREFIX is required")
	}

	rules := make([]envsplit.Rule, 0, len(ruleStrs))
	for _, rs := range ruleStrs {
		parts := strings.SplitN(rs, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid rule %q: expected bucket=PREFIX", rs)
		}
		rules = append(rules, envsplit.Rule{Bucket: parts[0], Prefix: parts[1]})
	}

	splitter, err := envsplit.New(rules)
	if err != nil {
		return err
	}

	var raw []byte
	if inputFile != "" {
		raw, err = os.ReadFile(inputFile)
	} else {
		raw, err = os.ReadFile("/dev/stdin")
	}
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	var secrets map[string]string
	if err := json.Unmarshal(raw, &secrets); err != nil {
		return fmt.Errorf("parsing secrets JSON: %w", err)
	}

	res := splitter.Apply(secrets)

	for _, name := range envsplit.BucketNames(res) {
		fmt.Printf("[%s]\n", name)
		for _, k := range sortedKeys(res.Buckets[name]) {
			fmt.Printf("  %s=%s\n", k, res.Buckets[name][k])
		}
	}

	if showUnrouted && len(res.Unrouted) > 0 {
		fmt.Println("[unrouted]")
		for _, k := range sortedKeys(res.Unrouted) {
			fmt.Printf("  %s=%s\n", k, res.Unrouted[k])
		}
	}
	return nil
}
