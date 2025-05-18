// internal/app/app.go
package app

import (
	"os"
	"path/filepath"
)

// Initialize sets up the application
func Initialize() error {
	// Initialize app state
	LogDebug("Initializing application")

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		LogDebug("Error getting current directory: %v", err)
		return err
	}

	// Set up environment variables for the device
	LogDebug("Setting environment variables")

	_ = os.Setenv("DEVICE", "brick")
	_ = os.Setenv("PLATFORM", "tg5040")

	// Add current directory to PATH instead of replacing it
	existingPath := os.Getenv("PATH")
	newPath := cwd + ":" + existingPath
	_ = os.Setenv("PATH", newPath)
	LogDebug("Updated PATH: %s", newPath)

	_ = os.Setenv("LD_LIBRARY_PATH", "/mnt/SDCARD/.system/tg5040/lib:/usr/trimui/lib")

	// Create required directories
	LogDebug("Creating application directories")

	// Create logs directory
	err = os.MkdirAll(filepath.Join(cwd, "Logs"), 0755)
	if err != nil {
		LogDebug("Error creating Logs directory: %v", err)
	}

	// Create themes directory
	err = os.MkdirAll(filepath.Join(cwd, "Themes"), 0755)
	if err != nil {
		LogDebug("Error creating Themes directory: %v", err)
	}

	// Create cache directory for temporary files
	err = os.MkdirAll(filepath.Join(cwd, ".cache"), 0755)
	if err != nil {
		LogDebug("Error creating .cache directory: %v", err)
	}

	LogDebug("Initialization complete")
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

// GetLogsDir returns the path to the logs directory
func GetLogsDir() string {
	return filepath.Join(GetWorkingDir(), "Logs")
}