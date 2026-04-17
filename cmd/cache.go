package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/cache"
)

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Manage the local secrets cache",
}

var cacheClearCmd = &cobra.Command{
	Use:   "clear [vault-path]",
	Short: "Invalidate the cache for a given Vault path",
	Args:  cobra.ExactArgs(1),
	RunE:  runCacheClear,
}

var cacheStatusCmd = &cobra.Command{
	Use:   "status [vault-path]",
	Short: "Show whether a cache entry exists and is valid",
	Args:  cobra.ExactArgs(1),
	RunE:  runCacheStatus,
}

var cacheTTL time.Duration
var cacheDir string

func init() {
	cacheCmd.PersistentFlags().StringVar(&cacheDir, "cache-dir", ".vaultline/cache", "Directory for cache files")
	cacheCmd.PersistentFlags().DurationVar(&cacheTTL, "ttl", 5*time.Minute, "Cache TTL duration")
	cacheCmd.AddCommand(cacheClearCmd)
	cacheCmd.AddCommand(cacheStatusCmd)
	rootCmd.AddCommand(cacheCmd)
}

func runCacheClear(cmd *cobra.Command, args []string) error {
	vaultPath := args[0]
	c := cache.New(cacheDir, cacheTTL)
	if err := c.Invalidate(vaultPath); err != nil {
		return fmt.Errorf("failed to invalidate cache: %w", err)
	}
	fmt.Fprintf(os.Stdout, "Cache cleared for path: %s\n", vaultPath)
	return nil
}

func runCacheStatus(cmd *cobra.Command, args []string) error {
	vaultPath := args[0]
	c := cache.New(cacheDir, cacheTTL)
	secrets, ok := c.Get(vaultPath)
	if !ok {
		fmt.Fprintf(os.Stdout, "Cache MISS for path: %s\n", vaultPath)
		return nil
	}
	fmt.Fprintf(os.Stdout, "Cache HIT for path: %s (%d keys cached)\n", vaultPath, len(secrets))
	return nil
}
