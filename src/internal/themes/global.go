// src/internal/themes/global.go
// Global theme operations

package themes

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"nextui-themes/internal/logging"
	"nextui-themes/internal/system"
)

// ListGlobalThemes returns a list of available global themes
func ListGlobalThemes(globalThemesDir string) ([]string, error) {
	var themes []string

	// Check if the directory exists
	_, err := os.Stat(globalThemesDir)
	if os.IsNotExist(err) {
		logging.LogDebug("Global themes directory does not exist: %s", globalThemesDir)
		return nil, fmt.Errorf("global themes directory does not exist: %s", globalThemesDir)
	} else if err != nil {
		logging.LogDebug("Error checking global themes directory: %v", err)
		return nil, fmt.Errorf("error checking themes directory: %w", err)
	}

	// Read the directory
	entries, err := os.ReadDir(globalThemesDir)
	if err != nil {
		logging.LogDebug("Error reading global themes directory: %v", err)
		return nil, fmt.Errorf("error reading themes directory: %w", err)
	}

	// Find directories that contain a bg.png file
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			bgPath := filepath.Join(globalThemesDir, entry.Name(), "bg.png")
			if _, err := os.Stat(bgPath); err == nil {
				themes = append(themes, entry.Name())
			}
		}
	}

	logging.LogDebug("Found %d global themes", len(themes))
	return themes, nil
}

// ApplyGlobalTheme applies a global theme to all directories
func ApplyGlobalTheme(themeName string) error {
	logging.LogDebug("Applying global theme: %s", themeName)

	// Get current directory for absolute paths
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}
	logging.LogDebug("Current directory: %s", cwd)

	// Get system paths
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		logging.LogDebug("Error getting system paths: %v", err)
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Ensure all media directories exist
	if err := system.EnsureMediaDirectories(systemPaths); err != nil {
		logging.LogDebug("Error ensuring media directories: %v", err)
		return fmt.Errorf("error ensuring media directories: %w", err)
	}

	// Source background image
	srcBg := filepath.Join(cwd, "Themes", "Global", themeName, "bg.png")
	logging.LogDebug("Theme background path: %s", srcBg)

	// Check if the source background exists
	_, err = os.Stat(srcBg)
	if err != nil {
		logging.LogDebug("Theme background not found: %v", err)
		return fmt.Errorf("theme background not found: %w", err)
	}

	// Apply to root .media directory
	rootMediaBg := filepath.Join(systemPaths.Root, ".media", "bg.png")
	logging.LogDebug("Copying to root .media: %s", rootMediaBg)
	if err := CopyFile(srcBg, rootMediaBg); err != nil {
		logging.LogDebug("Error copying to root .media: %v", err)
		return fmt.Errorf("failed to copy background to root .media: %w", err)
	}

	// Also apply to root directory (NextUI sometimes looks for bg.png in the root)
	rootBg := filepath.Join(systemPaths.Root, "bg.png")
	logging.LogDebug("Copying to root: %s", rootBg)
	if err := CopyFile(srcBg, rootBg); err != nil {
		logging.LogDebug("Error copying to root: %v", err)
		return fmt.Errorf("failed to copy background to root: %w", err)
	}

	// Apply to Recently Played
	rpBg := filepath.Join(systemPaths.RecentlyPlayed, ".media", "bg.png")
	logging.LogDebug("Copying to Recently Played: %s", rpBg)
	if err := CopyFile(srcBg, rpBg); err != nil {
		logging.LogDebug("Error copying to Recently Played: %v", err)
		return fmt.Errorf("failed to copy background to Recently Played: %w", err)
	}

	// Apply to Tools
	toolsBg := filepath.Join(systemPaths.Tools, ".media", "bg.png")
	logging.LogDebug("Copying to Tools: %s", toolsBg)
	if err := CopyFile(srcBg, toolsBg); err != nil {
		logging.LogDebug("Error copying to Tools: %v", err)
		return fmt.Errorf("failed to copy background to Tools: %w", err)
	}

	// Apply to all system directories
	for _, system := range systemPaths.Systems {
		systemBg := filepath.Join(system.MediaPath, "bg.png")
		logging.LogDebug("Copying to system %s: %s", system.Name, systemBg)
		if err := CopyFile(srcBg, systemBg); err != nil {
			logging.LogDebug("Error copying to system %s: %v", system.Name, err)
			return fmt.Errorf("failed to copy background to %s: %w", system.Name, err)
		}
	}

	logging.LogDebug("Global theme applied successfully")
	return nil
}

// ApplyCustomTheme applies a custom theme to a specific system
func ApplyCustomTheme(systemName string) error {
	logging.LogDebug("Applying custom theme to system: %s", systemName)

	// Get current directory for absolute paths
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Get system paths
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		logging.LogDebug("Error getting system paths: %v", err)
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Get all available backgrounds from Global themes
	globalThemesPath := filepath.Join(cwd, "Themes", "Global")
	logging.LogDebug("Scanning global themes directory: %s", globalThemesPath)

	// List themes is now its own function we can reuse
	themes, err := ListGlobalThemes(globalThemesPath)
	if err != nil {
		logging.LogDebug("Error reading global themes: %v", err)
		return fmt.Errorf("error reading global themes: %w", err)
	}

	if len(themes) == 0 {
		logging.LogDebug("No themes found in global themes directory")
		return fmt.Errorf("no themes found in Global directory")
	}

	// Instead of handling UI here, we'll have the UI layer handle theme selection
	// and just apply the selected theme by name

	// Since our function is called with a specific theme name, we can now apply it directly
	// This will be managed by the UI layer in theme_selection.go
	themeName := os.Getenv("SELECTED_THEME")
	if themeName == "" {
		logging.LogDebug("No theme selected for system: %s", systemName)
		return fmt.Errorf("no theme selected for system: %s", systemName)
	}

	// Source background
	srcBg := filepath.Join(globalThemesPath, themeName, "bg.png")
	logging.LogDebug("Selected theme: %s (%s)", themeName, srcBg)

	// Check if the source file exists
	if _, err := os.Stat(srcBg); err != nil {
		logging.LogDebug("Theme background file not found: %s, error: %v", srcBg, err)
		return fmt.Errorf("theme background file not found: %s", srcBg)
	}

	// Determine target directory
	var targetPath string
	var targetMediaPath string

	if systemName == "Root" {
		targetPath = systemPaths.Root
		targetMediaPath = filepath.Join(targetPath, ".media")
	} else if systemName == "Recently Played" {
		targetPath = systemPaths.RecentlyPlayed
		targetMediaPath = filepath.Join(targetPath, ".media")
	} else if systemName == "Tools" {
		targetPath = systemPaths.Tools
		targetMediaPath = filepath.Join(targetPath, ".media")
	} else {
		// Find the system in our list
		found := false
		for _, system := range systemPaths.Systems {
			if strings.Contains(system.Name, systemName) {
				targetPath = system.Path
				targetMediaPath = system.MediaPath
				found = true
				break
			}
		}

		if !found {
			logging.LogDebug("System not found: %s", systemName)
			return fmt.Errorf("system not found: %s", systemName)
		}
	}

	// Ensure media directory exists
	if err := os.MkdirAll(targetMediaPath, 0755); err != nil {
		logging.LogDebug("Error creating media directory: %v", err)
		return fmt.Errorf("failed to create media directory: %w", err)
	}

	// Apply background
	dstBg := filepath.Join(targetMediaPath, "bg.png")
	if err := CopyFile(srcBg, dstBg); err != nil {
		logging.LogDebug("Error copying background: %v", err)
		return fmt.Errorf("failed to copy background: %w", err)
	}

	// If applying to root, also copy to the root directory itself
	if systemName == "Root" {
		rootBg := filepath.Join(targetPath, "bg.png")
		if err := CopyFile(srcBg, rootBg); err != nil {
			logging.LogDebug("Error copying to root: %v", err)
			return fmt.Errorf("failed to copy background to root: %w", err)
		}
	}

	logging.LogDebug("Custom theme applied successfully for %s", systemName)
	return nil
}

// DisplayThemeSelectionList is a helper function for theme selection
func DisplayThemeSelectionList(themes []string, title string) (string, int) {
	// Implement using UI functions from common.go
	// This is a simple wrapper to avoid circular imports
	return displayMinUiList(strings.Join(themes, "\n"), "text", title)
}

// Placeholder implementation - will be replaced by actual UI function
func displayMinUiList(list string, format string, title string, extraArgs ...string) (string, int) {
	// This is a placeholder that would be replaced by the actual UI function
	// In a real implementation, this would either be provided by an interface
	// or we would need to restructure to avoid circular dependencies
	return "", 0
}