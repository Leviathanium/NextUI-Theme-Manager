// src/internal/app/app.go
// Simplified application initialization focused on theme management

package app

import (
	"os"
	"path/filepath"

	"nextui-themes/internal/logging"
	"nextui-themes/internal/themes"
)

// Initialize sets up the application
func Initialize() error {
	// Initialize app state
	state.CurrentScreen = MainMenu

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return err
	}

	// Set up environment variables for the TrimUI brick
	logging.LogDebug("Setting environment variables")

	_ = os.Setenv("DEVICE", "brick")
	_ = os.Setenv("PLATFORM", "tg5040")

	// Add current directory to PATH instead of replacing it
	existingPath := os.Getenv("PATH")
	newPath := cwd + ":" + existingPath
	_ = os.Setenv("PATH", newPath)
	logging.LogDebug("Updated PATH: %s", newPath)

	_ = os.Setenv("LD_LIBRARY_PATH", "/mnt/SDCARD/.system/tg5040/lib:/usr/trimui/lib")

	// Create theme directory structure
	logging.LogDebug("Creating theme directories")

	// Create Themes directory structure
	err = os.MkdirAll(filepath.Join(cwd, "Themes", "Imports"), 0755)
	if err != nil {
		logging.LogDebug("Error creating Themes/Imports directory: %v", err)
	}

	err = os.MkdirAll(filepath.Join(cwd, "Themes", "Exports"), 0755)
	if err != nil {
		logging.LogDebug("Error creating Themes/Exports directory: %v", err)
	}

	// Create logs directory
	err = os.MkdirAll(filepath.Join(cwd, "Logs"), 0755)
	if err != nil {
		logging.LogDebug("Error creating Logs directory: %v", err)
	}

	// Explicitly initialize theme directories after logging is set up
	if err := themes.EnsureThemeDirectoryStructure(); err != nil {
		logging.LogDebug("Warning: Could not create theme directories: %v", err)
	}

	if err := themes.CreatePlaceholderFiles(); err != nil {
		logging.LogDebug("Warning: Could not create placeholder files: %v", err)
	}

	// Log about theme functionality
	logging.LogDebug("Theme import/export functionality initialized")

	logging.LogDebug("Initialization complete")
	return nil
}