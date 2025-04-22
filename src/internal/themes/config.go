// src/internal/themes/config.go
// Configuration management for the theme manager

package themes

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"nextui-themes/internal/logging"
)

// ConfigData represents the configuration file structure
type ConfigData struct {
	RepoURL  string `json:"repo_url"`
	Branch   string `json:"branch"`
	Version  string `json:"version"`
	DeviceID string `json:"device_id,omitempty"`
}

// Default configuration values
const (
	DefaultRepoURL = "https://github.com/Leviathanium/NextUI-Themes"
	DefaultBranch  = "main"
	DefaultVersion = "1.0.0"
)

// LoadConfig loads the configuration from the config file
func LoadConfig() (*ConfigData, error) {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting current directory: %w", err)
	}

	// Path to config file
	configPath := filepath.Join(cwd, "config.json")

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config
		config := &ConfigData{
			RepoURL:  DefaultRepoURL,
			Branch:   DefaultBranch,
			Version:  DefaultVersion,
		}

		// Save default config
		if err := SaveConfig(config); err != nil {
			return nil, fmt.Errorf("error saving default config: %w", err)
		}

		return config, nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Parse config
	var config ConfigData
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	// Update repo settings
	if config.RepoURL != "" {
		SetRepoURL(config.RepoURL)
	}

	if config.Branch != "" {
		SetRepoBranch(config.Branch)
	}

	return &config, nil
}

// SaveConfig saves the configuration to the config file
func SaveConfig(config *ConfigData) error {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Path to config file
	configPath := filepath.Join(cwd, "config.json")

	// Convert to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	logging.LogDebug("Saved configuration to %s", configPath)
	return nil
}

// InitializeConfig loads the configuration and sets up the repo settings
func InitializeConfig() error {
	config, err := LoadConfig()
	if err != nil {
		logging.LogDebug("Error loading configuration: %v", err)
		logging.LogDebug("Using default repository settings")
	} else {
		logging.LogDebug("Loaded configuration: repo=%s, branch=%s", config.RepoURL, config.Branch)
	}

	return nil
}

// UpdateRepoURL updates the repository URL in the configuration
func UpdateRepoURL(url string) error {
	// Load current config
	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	// Update URL
	config.RepoURL = url
	SetRepoURL(url)

	// Save config
	return SaveConfig(config)
}

// UpdateRepoBranch updates the repository branch in the configuration
func UpdateRepoBranch(branch string) error {
	// Load current config
	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	// Update branch
	config.Branch = branch
	SetRepoBranch(branch)

	// Save config
	return SaveConfig(config)
}

func init() {
	// Initialize configuration when the package is loaded
	if err := InitializeConfig(); err != nil {
		logging.LogDebug("Warning: Failed to initialize configuration: %v", err)
	}
}