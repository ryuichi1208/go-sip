package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg == nil {
		t.Fatal("DefaultConfig() returned nil")
	}

	if cfg.Server.Port != "5060" {
		t.Errorf("Default port should be 5060, got %s", cfg.Server.Port)
	}

	if cfg.Server.LogLevel != "info" {
		t.Errorf("Default log level should be info, got %s", cfg.Server.LogLevel)
	}

	if cfg.Server.BindAddr != "0.0.0.0" {
		t.Errorf("Default bind address should be 0.0.0.0, got %s", cfg.Server.BindAddr)
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	// Ensure testdata directory exists
	testDir := filepath.Join("testdata")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create a temporary file path in testdata
	tempFile := filepath.Join(testDir, "test_config_save_load.json")
	defer os.Remove(tempFile)

	// Create and save a config
	originalCfg := &Config{
		Server: ServerConfig{
			Port:     "5070",
			LogLevel: "debug",
			BindAddr: "127.0.0.1",
		},
	}

	err := SaveConfig(originalCfg, tempFile)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load the config back
	loadedCfg, err := LoadConfig(tempFile)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Check values match
	if loadedCfg.Server.Port != originalCfg.Server.Port {
		t.Errorf("Port mismatch: %s != %s", loadedCfg.Server.Port, originalCfg.Server.Port)
	}

	if loadedCfg.Server.LogLevel != originalCfg.Server.LogLevel {
		t.Errorf("LogLevel mismatch: %s != %s", loadedCfg.Server.LogLevel, originalCfg.Server.LogLevel)
	}

	if loadedCfg.Server.BindAddr != originalCfg.Server.BindAddr {
		t.Errorf("BindAddr mismatch: %s != %s", loadedCfg.Server.BindAddr, originalCfg.Server.BindAddr)
	}
}

func TestLoadExistingConfig(t *testing.T) {
	// Use the test config file that should exist in testdata
	configFile := filepath.Join("testdata", "test_config.json")

	// If file doesn't exist in the test environment, create it
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Create testdata dir if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(configFile), 0755); err != nil {
			t.Skipf("Could not create testdata directory: %v", err)
		}

		// Create a sample config file for testing
		testCfg := &Config{
			Server: ServerConfig{
				Port:     "5070",
				LogLevel: "info",
				BindAddr: "127.0.0.1",
			},
		}

		if err := SaveConfig(testCfg, configFile); err != nil {
			t.Skipf("Could not create test config file: %v", err)
		}
	}

	// Load the config
	cfg, err := LoadConfig(configFile)
	if err != nil {
		t.Fatalf("Failed to load existing config: %v", err)
	}

	// Check that we got a valid config
	if cfg == nil {
		t.Fatal("LoadConfig() returned nil config")
	}
}

func TestLoadNonExistentConfig(t *testing.T) {
	// Attempt to load a non-existent file
	cfg, err := LoadConfig("testdata/non_existent_file.json")

	// Check that we get an error but still get a default config
	if err == nil {
		t.Error("Expected an error when loading non-existent file")
	}

	if cfg == nil {
		t.Fatal("LoadConfig() should return default config when file not found")
	}

	// Check that it's the default config
	defaultCfg := DefaultConfig()
	if cfg.Server.Port != defaultCfg.Server.Port {
		t.Errorf("Expected default port %s, got %s", defaultCfg.Server.Port, cfg.Server.Port)
	}
}
