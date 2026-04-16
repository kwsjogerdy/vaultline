package vault

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is a minimal Vault HTTP client.
type Client struct {
	Address string
	Token   string
	http    *http.Client
}

// NewClient creates a new Vault client.
func NewClient(address, token string) *Client {
	return &Client{
		Address: address,
		Token:   token,
		http:    &http.Client{Timeout: 10 * time.Second},
	}
}

// SecretData holds the key/value pairs returned from a KV secret path.
type SecretData map[string]string

// GetSecrets fetches secrets from a KV v2 mount path.
func (c *Client) GetSecrets(mountPath, secretPath string) (SecretData, error) {
	url := fmt.Sprintf("%s/v1/%s/data/%s", c.Address, mountPath, secretPath)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.Token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized: check your Vault token")
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("secret not found at path: %s/%s", mountPath, secretPath)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d from Vault", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	var result struct {
		Data struct {
			Data map[string]interface{} `json:"data"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	secrets := make(SecretData, len(result.Data.Data))
	for k, v := range result.Data.Data {
		secrets[k] = fmt.Sprintf("%v", v)
	}
	return secrets, nil
}
