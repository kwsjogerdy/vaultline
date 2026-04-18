// Package notify sends desktop or webhook notifications after a sync.
package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Config holds notification settings.
type Config struct {
	WebhookURL string
	Channel    string
	Timeout    time.Duration
}

// Notifier sends sync result notifications.
type Notifier struct {
	cfg    Config
	client *http.Client
}

// New creates a Notifier with the given config.
func New(cfg Config) *Notifier {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	return &Notifier{
		cfg:    cfg,
		client: &http.Client{Timeout: timeout},
	}
}

type payload struct {
	Channel string `json:"channel,omitempty"`
	Text    string `json:"text"`
}

// Send posts a notification message to the configured webhook.
func (n *Notifier) Send(message string) error {
	if n.cfg.WebhookURL == "" {
		return fmt.Errorf("notify: webhook URL is not configured")
	}
	p := payload{Channel: n.cfg.Channel, Text: message}
	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("notify: marshal payload: %w", err)
	}
	resp, err := n.client.Post(n.cfg.WebhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("notify: post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("notify: unexpected status %d", resp.StatusCode)
	}
	return nil
}

// SyncMessage builds a standard notification message for a sync result.
func SyncMessage(env string, added, changed, removed int, err error) string {
	if err != nil {
		return fmt.Sprintf("[vaultline] sync FAILED for %s: %v", env, err)
	}
	return fmt.Sprintf("[vaultline] sync OK for %s — +%d ~%d -%d keys", env, added, changed, removed)
}
