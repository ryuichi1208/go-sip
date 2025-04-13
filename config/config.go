package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config represents the SIP server configuration
type Config struct {
	Server ServerConfig `json:"server"`
}

// ServerConfig holds server-specific settings
type ServerConfig struct {
	Port     string `json:"port"`
	LogLevel string `json:"log_level"`
	BindAddr string `json:"bind_addr"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:     "5060",
			LogLevel: "info",
			BindAddr: "0.0.0.0",
		},
	}
}

// LoadConfig loads configuration from the specified path
func LoadConfig(path string) (*Config, error) {
	// Load default configuration
	cfg := DefaultConfig()

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, fmt.Errorf("config file not found: %s", path)
	}

	// Read file
	configFile, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("error reading config file: %v", err)
	}

	// Parse JSON
	if err := json.Unmarshal(configFile, cfg); err != nil {
		return cfg, fmt.Errorf("error parsing config file: %v", err)
	}

	return cfg, nil
}

// SaveConfig saves configuration to a file
func SaveConfig(cfg *Config, path string) error {
	// Convert to JSON
	configJSON, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("error encoding config to JSON: %v", err)
	}

	// Write to file
	if err := os.WriteFile(path, configJSON, 0644); err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}

	return nil
}
