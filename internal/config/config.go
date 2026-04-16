package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

// Config holds the vaultline configuration.
type Config struct {
	VaultAddr  string `mapstructure:"vault_addr"`
	VaultToken string `mapstructure:"vault_token"`
	SecretPath string `mapstructure:"secret_path"`
	EnvFile    string `mapstructure:"env_file"`
}

// Load reads configuration from a file and environment variables.
func Load(cfgFile string) (*Config, error) {
	v := viper.New()

	v.SetDefault("vault_addr", "http://127.0.0.1:8200")
	v.SetDefault("env_file", ".env")

	v.SetEnvPrefix("VAULTLINE")
	v.AutomaticEnv()

	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		v.SetConfigName(".vaultline")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath(os.Getenv("HOME"))
	}

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return nil, err
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if cfg.VaultAddr == "" {
		return nil, errors.New("vault_addr is required")
	}
	if cfg.VaultToken == "" {
		return nil, errors.New("vault_token is required (set VAULTLINE_VAULT_TOKEN or vault_token in config)")
	}
	if cfg.SecretPath == "" {
		return nil, errors.New("secret_path is required")
	}

	return &cfg, nil
}
