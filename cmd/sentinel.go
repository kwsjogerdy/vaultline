package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/envsentinel"
)

var (
	sentinelRulesFile string
)

func init() {
	sentinelCmd := &cobra.Command{
		Use:   "sentinel",
		Short: "Check secrets against sentinel alert rules",
		RunE:  runSentinel,
	}
	sentinelCmd.Flags().StringVarP(&sentinelRulesFile, "rules", "r", "", "JSON file containing sentinel rules (required)")
	_ = sentinelCmd.MarkFlagRequired("rules")
	rootCmd.AddCommand(sentinelCmd)
}

type sentinelRuleFile struct {
	Rules []envsentinel.Rule `json:"rules"`
}

func runSentinel(cmd *cobra.Command, args []string) error {
	data, err := os.ReadFile(sentinelRulesFile)
	if err != nil {
		return fmt.Errorf("reading rules file: %w", err)
	}

	var rf sentinelRuleFile
	if err := json.Unmarshal(data, &rf); err != nil {
		return fmt.Errorf("parsing rules file: %w", err)
	}

	s, err := envsentinel.New(rf.Rules)
	if err != nil {
		return fmt.Errorf("building sentinel: %w", err)
	}

	// Read secrets from environment (or pipe via stdin JSON).
	secrets := map[string]string{}
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		if err := json.NewDecoder(os.Stdin).Decode(&secrets); err != nil {
			return fmt.Errorf("decoding secrets from stdin: %w", err)
		}
	} else {
		for _, e := range os.Environ() {
			var k, v string
			for i, c := range e {
				if c == '=' {
					k = e[:i]
					v = e[i+1:]
					break
				}
			}
			if k != "" {
				secrets[k] = v
			}
		}
	}

	violations := s.Check(secrets)
	if len(violations) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "sentinel: all checks passed")
		return nil
	}

	for _, v := range violations {
		fmt.Fprintln(cmd.OutOrStdout(), v.String())
	}
	return fmt.Errorf("sentinel: %d violation(s) found", len(violations))
}
