// internal/app/app.go
package app

import (
	"os"
	"path/filepath"
	"encoding/json"

	"thememanager/internal/logging"
)

// Settings structure for application configuration
type Settings struct {
	AutoBackup bool   `json:"auto_backup"`
	Version    string `json:"version"`
}

// Global settings
var appSettings Settings

// Initialize sets up the application
func Initialize() error {
	// Initialize app state
	logging.LogDebug("Initializing application")

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return err
	}

	// Set up environment variables for the device
	logging.LogDebug("Setting environment variables")

	_ = os.Setenv("DEVICE", "brick")
	_ = os.Setenv("PLATFORM", "tg5040")

	// Add current directory to PATH instead of replacing it
	existingPath := os.Getenv("PATH")
	newPath := cwd + ":" + existingPath
	_ = os.Setenv("PATH", newPath)
	logging.LogDebug("Updated PATH: %s", newPath)

	_ = os.Setenv("LD_LIBRARY_PATH", "/mnt/SDCARD/.system/tg5040/lib:/usr/trimui/lib")

	// Create required directories
	logging.LogDebug("Creating application directories")

	// Create logs directory
	err = os.MkdirAll(filepath.Join(cwd, "Logs"), 0755)
	if err != nil {
		logging.LogDebug("Error creating Logs directory: %v", err)
	}

	// Create cache directory for temporary files
	err = os.MkdirAll(filepath.Join(cwd, ".cache"), 0755)
	if err != nil {
		logging.LogDebug("Error creating .cache directory: %v", err)
	}

	// Load settings
	loadSettings()

	logging.LogDebug("Initialization complete")
	return nil
}

// GetWorkingDir returns the current working directory
func GetWorkingDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return cwd
}

// GetThemesDir returns the path to the themes directory
func GetThemesDir() string {
	return filepath.Join(GetWorkingDir(), "Themes")
}

// GetOverlaysDir returns the path to the overlays directory
func GetOverlaysDir() string {
	return filepath.Join(GetWorkingDir(), "Overlays")
}

// GetBackupsDir returns the path to the backups directory
func GetBackupsDir() string {
	return filepath.Join(GetWorkingDir(), "Backups")
}

// GetCatalogDir returns the path to the catalog directory
func GetCatalogDir() string {
	return filepath.Join(GetWorkingDir(), "Catalog")
}

// GetAutoBackupSetting returns the auto-backup setting
func GetAutoBackupSetting() bool {
	return appSettings.AutoBackup
}

// SetAutoBackupSetting updates the auto-backup setting
func SetAutoBackupSetting(enabled bool) {
	appSettings.AutoBackup = enabled
	saveSettings()
}

// loadSettings loads application settings from file
func loadSettings() {
	// Default settings
	appSettings = Settings{
		AutoBackup: false,
		Version:    "1.0.0",
	}

	// Get settings file path
	settingsPath := filepath.Join(GetWorkingDir(), "settings.json")

	// Check if file exists
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		// Create default settings file
		saveSettings()
		return
	}

	// Read settings file
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		logging.LogDebug("Error reading settings file: %v", err)
		return
	}

	// Parse settings
	if err := json.Unmarshal(data, &appSettings); err != nil {
		logging.LogDebug("Error parsing settings: %v", err)
		return
	}

	logging.LogDebug("Settings loaded: auto-backup=%v", appSettings.AutoBackup)
}

// saveSettings saves application settings to file
func saveSettings() {
	// Get settings file path
	settingsPath := filepath.Join(GetWorkingDir(), "settings.json")

	// Convert settings to JSON
	data, err := json.MarshalIndent(appSettings, "", "  ")
	if err != nil {
		logging.LogDebug("Error converting settings to JSON: %v", err)
		return
	}

	// Write settings file
	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		logging.LogDebug("Error writing settings file: %v", err)
		return
	}

	logging.LogDebug("Settings saved: auto-backup=%v", appSettings.AutoBackup)
}