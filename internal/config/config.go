package config

import (
	"errors"
	"os"

	"github.com/vaultline/vaultline/internal/filter"
	"gopkg.in/yaml.v3"
)

// Config holds all configuration for a vaultline run.
type Config struct {
	VaultAddr  string   `yaml:"vault_addr"`
	VaultToken string   `yaml:"vault_token"`
	SecretPath string   `yaml:"secret_path"`
	OutputFile string   `yaml:"output_file"`
	Prefixes   []string `yaml:"prefixes"`
	Excludes   []string `yaml:"excludes"`
}

// FilterRule converts the config's prefix/exclude lists into a filter.Rule.
func (c *Config) FilterRule() filter.Rule {
	return filter.Rule{
		Prefixes: c.Prefixes,
		Excludes: c.Excludes,
	}
}

// Load reads config from a YAML file, then overrides with environment variables.
func Load(path string) (*Config, error) {
	cfg := &Config{}

	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		if err == nil {
			if err := yaml.Unmarshal(data, cfg); err != nil {
				return nil, err
			}
		}
	}

	if v := os.Getenv("VAULT_ADDR"); v != "" {
		cfg.VaultAddr = v
	}
	if v := os.Getenv("VAULT_TOKEN"); v != "" {
		cfg.VaultToken = v
	}
	if v := os.Getenv("VAULT_SECRET_PATH"); v != "" {
		cfg.SecretPath = v
	}
	if v := os.Getenv("VAULTLINE_OUTPUT"); v != "" {
		cfg.OutputFile = v
	}

	if cfg.VaultToken == "" {
		return nil, errors.New("vault token is required (set VAULT_TOKEN or vault_token in config)")
	}
	if cfg.SecretPath == "" {
		return nil, errors.New("secret path is required")
	}

	if cfg.OutputFile == "" {
		cfg.OutputFile = ".env"
	}

	return cfg, nil
}
