// Package envimport provides functionality for importing secrets from
// existing .env files into Vault.
package envimport

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/vaultline/vaultline/internal/env"
)

// Importer reads a .env file and writes its key/value pairs to Vault.
type Importer struct {
	vaultAddr  string
	vaultToken string
	client     *http.Client
}

// New creates a new Importer.
func New(vaultAddr, vaultToken string) *Importer {
	return &Importer{
		vaultAddr:  strings.TrimRight(vaultAddr, "/"),
		vaultToken: vaultToken,
		client:     &http.Client{},
	}
}

// Result holds the outcome of an import operation.
type Result struct {
	Imported []string
	Skipped  []string
	Errors   []string
}

// FromFile reads the given .env file and imports all secrets to the
// specified Vault path. Keys listed in skip are ignored.
func (im *Importer) FromFile(filePath, vaultPath string, skip []string) (*Result, error) {
	secrets, err := env.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("envimport: read file: %w", err)
	}

	skipSet := make(map[string]struct{}, len(skip))
	for _, k := range skip {
		skipSet[k] = struct{}{}
	}

	result := &Result{}
	for k, v := range secrets {
		if _, ok := skipSet[k]; ok {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		if err := im.writeSecret(vaultPath, k, v); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", k, err))
			continue
		}
		result.Imported = append(result.Imported, k)
	}
	return result, nil
}

func (im *Importer) writeSecret(vaultPath, key, value string) error {
	url := fmt.Sprintf("%s/v1/%s", im.vaultAddr, strings.Trim(vaultPath, "/"))
	body := fmt.Sprintf(`{"data":{"%s":"%s"}}`, key, value)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("X-Vault-Token", im.vaultToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := im.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("vault returned status %d", resp.StatusCode)
	}
	return nil
}
