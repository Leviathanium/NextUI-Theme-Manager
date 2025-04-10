// src/internal/app/app.go
// Application initialization and setup

package app

import (
	"os"
	"path/filepath"

	"nextui-themes/internal/logging"
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

	// Create theme directories if they don't exist
	logging.LogDebug("Creating theme directories")

	err = os.MkdirAll(filepath.Join(cwd, "Themes", "Global"), 0755)
	if err != nil {
		logging.LogDebug("Error creating Global themes directory: %v", err)
	}

	err = os.MkdirAll(filepath.Join(cwd, "Themes", "Dynamic"), 0755)
	if err != nil {
		logging.LogDebug("Error creating Dynamic themes directory: %v", err)
	}

	err = os.MkdirAll(filepath.Join(cwd, "Themes", "Default"), 0755)
	if err != nil {
		logging.LogDebug("Error creating Default themes directory: %v", err)
	}

	logging.LogDebug("Initialization complete")
	return nil
}