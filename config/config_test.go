package config

import (
	"os"
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
	// Create a temporary file path
	tempFile := "test_config.json"
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

func TestLoadNonExistentConfig(t *testing.T) {
	// Attempt to load a non-existent file
	cfg, err := LoadConfig("non_existent_file.json")

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
