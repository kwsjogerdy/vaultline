package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultline/vaultline/internal/notify"
)

var (
	notifyWebhook string
	notifyChannel string
	notifyMessage string
)

func init() {
	notifyCmd := &cobra.Command{
		Use:   "notify",
		Short: "Send a notification to a configured webhook",
		RunE:  runNotify,
	}
	notifyCmd.Flags().StringVar(&notifyWebhook, "webhook", os.Getenv("VAULTLINE_WEBHOOK_URL"), "Webhook URL (or VAULTLINE_WEBHOOK_URL)")
	notifyCmd.Flags().StringVar(&notifyChannel, "channel", "", "Optional channel name")
	notifyCmd.Flags().StringVar(&notifyMessage, "message", "", "Message to send (required)")
	_ = notifyCmd.MarkFlagRequired("message")
	rootCmd.AddCommand(notifyCmd)
}

func runNotify(cmd *cobra.Command, args []string) error {
	if notifyWebhook == "" {
		return fmt.Errorf("webhook URL is required (--webhook or VAULTLINE_WEBHOOK_URL)")
	}
	n := notify.New(notify.Config{
		WebhookURL: notifyWebhook,
		Channel:    notifyChannel,
	})
	if err := n.Send(notifyMessage); err != nil {
		return fmt.Errorf("notify: %w", err)
	}
	fmt.Fprintln(cmd.OutOrStdout(), "notification sent")
	return nil
}
