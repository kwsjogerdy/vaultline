package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/vaultline/vaultline/internal/ratelimit"
)

var (
	rlRate  float64
	rlBurst int
	rlCalls int
)

func init() {
	rlCmd := &cobra.Command{
		Use:   "ratelimit",
		Short: "Test or demonstrate the rate limiter behaviour",
		RunE:  runRatelimit,
	}
	rlCmd.Flags().Float64VarP(&rlRate, "rate", "r", 2.0, "Tokens replenished per second")
	rlCmd.Flags().IntVarP(&rlBurst, "burst", "b", 5, "Maximum burst size")
	rlCmd.Flags().IntVarP(&rlCalls, "calls", "n", 10, "Number of simulated calls")
	rootCmd.AddCommand(rlCmd)
}

func runRatelimit(cmd *cobra.Command, args []string) error {
	limiter := ratelimit.New(rlRate, rlBurst)

	allowed := 0
	denied := 0

	for i := 1; i <= rlCalls; i++ {
		if err := limiter.Allow(); err != nil {
			denied++
			fmt.Printf("call %s: DENIED — %v\n", pad(i), err)
		} else {
			allowed++
			fmt.Printf("call %s: ALLOWED\n", pad(i))
		}
	}

	fmt.Printf("\nSummary: %d allowed, %d denied out of %d calls\n", allowed, denied, rlCalls)
	return nil
}

func pad(n int) string {
	s := strconv.Itoa(n)
	if len(s) < 2 {
		return " " + s
	}
	return s
}
