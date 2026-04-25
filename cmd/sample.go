package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/env"
	"github.com/vaultline/vaultline/internal/envsampler"
)

var (
	sampleN      int
	sampleSeed   int64
	samplePrefix string
	sampleFile   string
)

func init() {
	sampleCmd := &cobra.Command{
		Use:   "sample",
		Short: "Print a random sample of keys from a .env file",
		RunE:  runSample,
	}
	sampleCmd.Flags().IntVarP(&sampleN, "count", "n", 5, "Number of keys to sample")
	sampleCmd.Flags().Int64Var(&sampleSeed, "seed", 0, "Random seed (0 = non-deterministic)")
	sampleCmd.Flags().StringVar(&samplePrefix, "prefix", "", "Only sample keys with this prefix")
	sampleCmd.Flags().StringVarP(&sampleFile, "file", "f", ".env", "Source .env file")
	rootCmd.AddCommand(sampleCmd)
}

func runSample(cmd *cobra.Command, _ []string) error {
	secrets, err := env.ReadFile(sampleFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", sampleFile, err)
		return err
	}

	s, err := envsampler.New(envsampler.Options{
		N:      sampleN,
		Seed:   sampleSeed,
		Prefix: samplePrefix,
	})
	if err != nil {
		return err
	}

	sampled := s.Sample(secrets)
	keys := make([]string, 0, len(sampled))
	for k := range sampled {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, sampled[k])
	}
	return nil
}
