package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	ConfigVersion = 1
	DefaultAPIURL = "https://ygm.app"
)

// Config represents the CLI configuration stored in ~/.config/ygm/config.yml
type Config struct {
	Version    int                `yaml:"version"`
	DefaultOrg string             `yaml:"default_org,omitempty"`
	APIURL     string             `yaml:"api_url"`
	Accounts   map[string]Account `yaml:"accounts"`
}

// Account represents a logged-in organization
type Account struct {
	Token     string `yaml:"token"`
	UserEmail string `yaml:"user_email"`
	OrgID     int    `yaml:"org_id"`
	OrgName   string `yaml:"org_name"`
}

// ConfigPath returns the path to the config file
func ConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		// Fall back to ~/.config
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("could not determine config directory: %w", err)
		}
		configDir = filepath.Join(homeDir, ".config")
	}

	return filepath.Join(configDir, "ygm", "config.yml"), nil
}

// Load reads the config from disk
func Load() (*Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No config yet
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

// Save writes the config to disk with secure permissions
func (c *Config) Save() error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}

	// Write with secure permissions (0600 = owner read/write only)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// AddAccount adds or updates an account in the config
func (c *Config) AddAccount(slug string, account Account) {
	if c.Accounts == nil {
		c.Accounts = make(map[string]Account)
	}
	c.Accounts[slug] = account

	// Set as default if it's the first account
	if c.DefaultOrg == "" {
		c.DefaultOrg = slug
	}
}

// NewConfig creates a new config with defaults
func NewConfig() *Config {
	return &Config{
		Version:  ConfigVersion,
		APIURL:   DefaultAPIURL,
		Accounts: make(map[string]Account),
	}
}
