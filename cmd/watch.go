package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourorg/vaultline/internal/watch"
)

var (
	watchInterval stringtRunE:  runWatch,
	}
	watchCmd.Flags().StringVar(&watchInterval, "interval", "5s", "Poll interval (e.g. 5s, 1m)")
	RootCmd.AddCommand(watchCmd)
}

func runWatch(cmd *cobra.Command, args []string) error {
	interval, err := time.ParseDuration(watchInterval)
	if err != nil {
		return fmt.Errorf("invalid interval %q: %w", watchInterval, err)
	}

	cfgPath := os.Getenv("VAULTLINE_CONFIG")
	if cfgPath == "" {
		cfgPath = ".vaultline.yaml"
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Watching %s every %s...\n", cfgPath, interval)

	w := watch.New(cfgPath, interval, func() error {
		fmt.Fprintln(cmd.OutOrStdout(), "Change detected — running sync")
		return runSync(cmd, args)
	})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := w.Run(ctx); err != nil && err != context.Canceled {
		return err
	}
	fmt.Fprintln(cmd.OutOrStdout(), "Watch stopped.")
	return nil
}
