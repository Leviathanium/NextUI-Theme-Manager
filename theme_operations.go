package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	// Background image name
	bgImageName = "bg.png"

	// Themes directories - CORRECTED!
	themesDir      = "Themes"
	globalThemesDir = "Global"    // Changed from staticThemesDir/Static
	dynamicThemesDir = "Dynamic"  // Added for clarity
	defaultThemesDir = "Default"  // Changed from systemThemesDir/System
)

// Update ApplyStaticTheme with better path handling and logging
// This applies a Global theme (renamed from Static)
func ApplyStaticTheme(themeName string) error {
	LogDebug("Applying global theme: %s", themeName)

	// Get current directory for absolute paths
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}
	LogDebug("Current directory: %s", cwd)

	// Get system paths
	systemPaths, err := GetSystemPaths()
	if err != nil {
		LogDebug("Error getting system paths: %v", err)
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Ensure all media directories exist
	if err := EnsureMediaDirectories(systemPaths); err != nil {
		LogDebug("Error ensuring media directories: %v", err)
		return fmt.Errorf("error ensuring media directories: %w", err)
	}

	// Source background image - CORRECTED PATH
	srcBg := filepath.Join(cwd, themesDir, globalThemesDir, themeName, bgImageName)
	LogDebug("Theme background path: %s", srcBg)

	// Check if the source background exists
	_, err = os.Stat(srcBg)
	if err != nil {
		LogDebug("Theme background not found: %v", err)
		return fmt.Errorf("theme background not found: %w", err)
	}

	// Apply to root .media directory
	rootMediaBg := filepath.Join(systemPaths.Root, ".media", bgImageName)
	LogDebug("Copying to root .media: %s", rootMediaBg)
	if err := CopyFile(srcBg, rootMediaBg); err != nil {
		LogDebug("Error copying to root .media: %v", err)
		return fmt.Errorf("failed to copy background to root .media: %w", err)
	}

	// Also apply to root directory (NextUI sometimes looks for bg.png in the root)
	rootBg := filepath.Join(systemPaths.Root, bgImageName)
	LogDebug("Copying to root: %s", rootBg)
	if err := CopyFile(srcBg, rootBg); err != nil {
		LogDebug("Error copying to root: %v", err)
		return fmt.Errorf("failed to copy background to root: %w", err)
	}

	// Apply to Recently Played
	rpBg := filepath.Join(systemPaths.RecentlyPlayed, ".media", bgImageName)
	LogDebug("Copying to Recently Played: %s", rpBg)
	if err := CopyFile(srcBg, rpBg); err != nil {
		LogDebug("Error copying to Recently Played: %v", err)
		return fmt.Errorf("failed to copy background to Recently Played: %w", err)
	}

	// Apply to Tools
	toolsBg := filepath.Join(systemPaths.Tools, ".media", bgImageName)
	LogDebug("Copying to Tools: %s", toolsBg)
	if err := CopyFile(srcBg, toolsBg); err != nil {
		LogDebug("Error copying to Tools: %v", err)
		return fmt.Errorf("failed to copy background to Tools: %w", err)
	}

	// Apply to all system directories
	for _, system := range systemPaths.Systems {
		systemBg := filepath.Join(system.MediaPath, bgImageName)
		LogDebug("Copying to system %s: %s", system.Name, systemBg)
		if err := CopyFile(srcBg, systemBg); err != nil {
			LogDebug("Error copying to system %s: %v", system.Name, err)
			return fmt.Errorf("failed to copy background to %s: %w", system.Name, err)
		}
	}

	LogDebug("Global theme applied successfully")
	return nil
}

// Update RemoveAllBackgrounds with better logging
// This applies the Default black theme
func RemoveAllBackgrounds() error {
	LogDebug("Applying default black theme")

	// Get current directory for absolute paths
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Default black background path
	defaultBg := filepath.Join(cwd, themesDir, defaultThemesDir, "bg.png")
	LogDebug("Default theme background path: %s", defaultBg)

	// Check if the default background exists
	_, err = os.Stat(defaultBg)
	if err != nil {
		LogDebug("Default background not found: %v", err)
		return fmt.Errorf("default background not found: %w", err)
	}

	// Get system paths
	systemPaths, err := GetSystemPaths()
	if err != nil {
		LogDebug("Error getting system paths: %v", err)
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Apply to root .media directory
	rootMediaBg := filepath.Join(systemPaths.Root, ".media", bgImageName)
	LogDebug("Copying to root .media: %s", rootMediaBg)
	if err := CopyFile(defaultBg, rootMediaBg); err != nil {
		LogDebug("Error copying to root .media: %v", err)
		return fmt.Errorf("failed to copy background to root .media: %w", err)
	}

	// Also apply to root directory
	rootBg := filepath.Join(systemPaths.Root, bgImageName)
	LogDebug("Copying to root: %s", rootBg)
	if err := CopyFile(defaultBg, rootBg); err != nil {
		LogDebug("Error copying to root: %v", err)
		return fmt.Errorf("failed to copy background to root: %w", err)
	}

	// Apply to Recently Played
	rpBg := filepath.Join(systemPaths.RecentlyPlayed, ".media", bgImageName)
	LogDebug("Copying to Recently Played: %s", rpBg)
	if err := CopyFile(defaultBg, rpBg); err != nil {
		LogDebug("Error copying to Recently Played: %v", err)
		return fmt.Errorf("failed to copy background to Recently Played: %w", err)
	}

	// Apply to Tools
	toolsBg := filepath.Join(systemPaths.Tools, ".media", bgImageName)
	LogDebug("Copying to Tools: %s", toolsBg)
	if err := CopyFile(defaultBg, toolsBg); err != nil {
		LogDebug("Error copying to Tools: %v", err)
		return fmt.Errorf("failed to copy background to Tools: %w", err)
	}

	// Apply to all system directories
	for _, system := range systemPaths.Systems {
		systemBg := filepath.Join(system.MediaPath, bgImageName)
		LogDebug("Copying to system %s: %s", system.Name, systemBg)
		if err := CopyFile(defaultBg, systemBg); err != nil {
			LogDebug("Error copying to system %s: %v", system.Name, err)
			return fmt.Errorf("failed to copy background to %s: %w", system.Name, err)
		}
	}

	LogDebug("Default theme applied successfully")
	return nil
}

// Update CopyFile with better logging
func CopyFile(src, dst string) error {
	LogDebug("Copying %s to %s", src, dst)

	// Create the destination directory if it doesn't exist
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		LogDebug("Error creating directory %s: %v", dstDir, err)
		return fmt.Errorf("failed to create directory %s: %w", dstDir, err)
	}

	srcFile, err := os.Open(src)
	if err != nil {
		LogDebug("Error opening source file: %v", err)
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		LogDebug("Error creating destination file: %v", err)
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		LogDebug("Error copying file: %v", err)
		return fmt.Errorf("failed to copy file: %w", err)
	}

	LogDebug("Successfully copied %d bytes", bytes)
	return nil
}

// Applying a system-specific theme
// Renamed from ApplySystemTheme to CustomThemeSelection to match the intended functionality
func CustomThemeSelection(systemName string) error {
	LogDebug("Selecting custom theme for system: %s", systemName)

	// Get current directory for absolute paths
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Get system paths
	systemPaths, err := GetSystemPaths()
	if err != nil {
		LogDebug("Error getting system paths: %v", err)
		return fmt.Errorf("error getting system paths: %w", err)
	}

	// Get all available backgrounds from Global themes
	globalThemesPath := filepath.Join(cwd, themesDir, globalThemesDir)
	LogDebug("Scanning global themes directory: %s", globalThemesPath)

	entries, err := os.ReadDir(globalThemesPath)
	if err != nil {
		LogDebug("Error reading global themes directory: %v", err)
		return fmt.Errorf("error reading global themes directory: %w", err)
	}

	// Extract theme folder names
	var themes []string
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			bgPath := filepath.Join(globalThemesPath, entry.Name(), bgImageName)
			if _, err := os.Stat(bgPath); err == nil {
				themes = append(themes, entry.Name())
			}
		}
	}

	if len(themes) == 0 {
		LogDebug("No themes found in global themes directory")
		return fmt.Errorf("no themes found in Global directory")
	}

	// Display theme selection menu
	LogDebug("Displaying theme selection menu for %s", systemName)
	selection, exitCode := displayMinUiList(
		strings.Join(themes, "\n"),
		"text",
		fmt.Sprintf("Select Theme for %s", systemName),
	)

	if exitCode != 0 {
		LogDebug("Theme selection cancelled")
		return fmt.Errorf("theme selection cancelled")
	}

	// Source background
	srcBg := filepath.Join(globalThemesPath, selection, bgImageName)
	LogDebug("Selected theme: %s (%s)", selection, srcBg)

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
			LogDebug("System not found: %s", systemName)
			return fmt.Errorf("system not found: %s", systemName)
		}
	}

	// Ensure media directory exists
	if err := os.MkdirAll(targetMediaPath, 0755); err != nil {
		LogDebug("Error creating media directory: %v", err)
		return fmt.Errorf("failed to create media directory: %w", err)
	}

	// Apply background
	dstBg := filepath.Join(targetMediaPath, bgImageName)
	if err := CopyFile(srcBg, dstBg); err != nil {
		LogDebug("Error copying background: %v", err)
		return fmt.Errorf("failed to copy background: %w", err)
	}

	// If applying to root, also copy to the root directory itself
	if systemName == "Root" {
		rootBg := filepath.Join(targetPath, bgImageName)
		if err := CopyFile(srcBg, rootBg); err != nil {
			LogDebug("Error copying to root: %v", err)
			return fmt.Errorf("failed to copy background to root: %w", err)
		}
	}

	LogDebug("Custom theme applied successfully for %s", systemName)
	return nil
}