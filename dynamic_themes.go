package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ThemeFile represents a background image in a dynamic theme
type ThemeFile struct {
	SourcePath string // Path to the source image
	TargetPath string // Path where it should be copied
}

// Update ApplyDynamicTheme with better path handling and logging
func ApplyDynamicTheme(themeName string) error {
	LogDebug("Applying dynamic theme: %s", themeName)

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

	// Scan the dynamic theme
	themeFiles, err := ScanDynamicTheme(themeName, systemPaths)
	if err != nil {
		LogDebug("Error scanning theme: %v", err)
		return fmt.Errorf("error scanning theme: %w", err)
	}

	LogDebug("Found %d theme files to apply", len(themeFiles))

	// Apply each theme file
	for _, file := range themeFiles {
		LogDebug("Copying %s to %s", file.SourcePath, file.TargetPath)

		// Create parent directories if needed
		parentDir := filepath.Dir(file.TargetPath)
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			LogDebug("Error creating directory %s: %v", parentDir, err)
			return fmt.Errorf("error creating directory %s: %w", parentDir, err)
		}

		// Copy the file
		if err := CopyFile(file.SourcePath, file.TargetPath); err != nil {
			LogDebug("Error copying %s to %s: %v", file.SourcePath, file.TargetPath, err)
			return fmt.Errorf("error copying %s to %s: %w", file.SourcePath, file.TargetPath, err)
		}
	}

	LogDebug("Dynamic theme applied successfully")
	return nil
}

// Update ScanDynamicTheme with better path handling and logging
func ScanDynamicTheme(themeName string, systemPaths *SystemPaths) ([]ThemeFile, error) {
	var themeFiles []ThemeFile

	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting current directory: %w", err)
	}

	// Get the theme base directory (using absolute path)
	// CORRECTED: Changed the directory structure to match intended layout
	themeDir := filepath.Join(cwd, "Themes", "Dynamic", themeName)
	LogDebug("Scanning dynamic theme directory: %s", themeDir)

	// Check if the theme directory exists
	_, err = os.Stat(themeDir)
	if os.IsNotExist(err) {
		LogDebug("Theme directory does not exist: %s", themeDir)
		return nil, fmt.Errorf("theme directory does not exist: %s", themeDir)
	} else if err != nil {
		LogDebug("Error checking theme directory: %v", err)
		return nil, fmt.Errorf("error checking theme directory: %w", err)
	}

	// Handle Root background
	rootBg := filepath.Join(themeDir, "Root", "bg.png")
	if _, err := os.Stat(rootBg); err == nil {
		LogDebug("Found Root background: %s", rootBg)
		// Add root background to theme files
		themeFiles = append(themeFiles, ThemeFile{
			SourcePath: rootBg,
			TargetPath: filepath.Join(systemPaths.Root, ".media", "bg.png"),
		})

		// Also add to root directory (NextUI sometimes looks for bg.png in the root)
		themeFiles = append(themeFiles, ThemeFile{
			SourcePath: rootBg,
			TargetPath: filepath.Join(systemPaths.Root, "bg.png"),
		})
	} else {
		LogDebug("Root background not found: %s, error: %v", rootBg, err)
	}

	// Handle Recently Played background
	rpBg := filepath.Join(themeDir, "Recently Played", "bg.png")
	if _, err := os.Stat(rpBg); err == nil {
		LogDebug("Found Recently Played background: %s", rpBg)
		// Add Recently Played background to theme files
		themeFiles = append(themeFiles, ThemeFile{
			SourcePath: rpBg,
			TargetPath: filepath.Join(systemPaths.RecentlyPlayed, ".media", "bg.png"),
		})
	} else {
		LogDebug("Recently Played background not found: %s, error: %v", rpBg, err)
	}

	// Handle Tools background
	toolsBg := filepath.Join(themeDir, "Tools", "bg.png")
	if _, err := os.Stat(toolsBg); err == nil {
		LogDebug("Found Tools background: %s", toolsBg)
		// Add Tools background to theme files
		themeFiles = append(themeFiles, ThemeFile{
			SourcePath: toolsBg,
			TargetPath: filepath.Join(systemPaths.Tools, ".media", "bg.png"),
		})
	} else {
		LogDebug("Tools background not found: %s, error: %v", toolsBg, err)
	}

	// Handle ROM systems backgrounds
	romsDir := filepath.Join(themeDir, "Roms")
	if _, err := os.Stat(romsDir); err == nil {
		LogDebug("Found Roms directory: %s", romsDir)

		// Look for a default background for systems
		defaultBg := filepath.Join(romsDir, "default.png")
		hasDefaultBg := false
		if _, err := os.Stat(defaultBg); err == nil {
			LogDebug("Found default system background: %s", defaultBg)
			hasDefaultBg = true
		} else {
			LogDebug("Default system background not found: %s, error: %v", defaultBg, err)
		}

		// Iterate through each installed system
		for _, system := range systemPaths.Systems {
			LogDebug("Processing system: %s (tag: %s)", system.Name, system.Tag)
			foundBg := false

			// Try to find a matching background for this system
			// First try by tag (preferred)
			if system.Tag != "" {
				tagBg := filepath.Join(romsDir, system.Tag, "bg.png")
				if _, err := os.Stat(tagBg); err == nil {
					LogDebug("Found system background by tag: %s", tagBg)
					themeFiles = append(themeFiles, ThemeFile{
						SourcePath: tagBg,
						TargetPath: filepath.Join(system.MediaPath, "bg.png"),
					})
					foundBg = true
				} else {
					LogDebug("System background by tag not found: %s, error: %v", tagBg, err)
				}
			}

			// If not found by tag, try by full system name
			if !foundBg {
				nameBg := filepath.Join(romsDir, system.Name, "bg.png")
				if _, err := os.Stat(nameBg); err == nil {
					LogDebug("Found system background by name: %s", nameBg)
					themeFiles = append(themeFiles, ThemeFile{
						SourcePath: nameBg,
						TargetPath: filepath.Join(system.MediaPath, "bg.png"),
					})
					foundBg = true
				} else {
					LogDebug("System background by name not found: %s, error: %v", nameBg, err)
				}
			}

			// If still not found, use the default background if available
			if !foundBg && hasDefaultBg {
				LogDebug("Using default background for system: %s", system.Name)
				themeFiles = append(themeFiles, ThemeFile{
					SourcePath: defaultBg,
					TargetPath: filepath.Join(system.MediaPath, "bg.png"),
				})
			}
		}
	} else {
		LogDebug("Roms directory not found: %s, error: %v", romsDir, err)
	}

	LogDebug("Theme scan complete, found %d files", len(themeFiles))
	return themeFiles, nil
}

// Update ListDynamicThemes with better path handling and logging
func ListDynamicThemes() ([]string, error) {
	var themes []string

	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		LogDebug("Error getting current directory: %v", err)
		return nil, fmt.Errorf("error getting current directory: %w", err)
	}

	// Get the dynamic themes directory (using absolute path)
	themesDir := filepath.Join(cwd, "Themes", "Dynamic")
	LogDebug("Listing themes from directory: %s", themesDir)

	// Check if the directory exists
	_, err = os.Stat(themesDir)
	if os.IsNotExist(err) {
		LogDebug("Dynamic themes directory does not exist: %s", themesDir)
		return nil, fmt.Errorf("dynamic themes directory does not exist: %s", themesDir)
	} else if err != nil {
		LogDebug("Error checking dynamic themes directory: %v", err)
		return nil, fmt.Errorf("error checking themes directory: %w", err)
	}

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(themesDir, 0755); err != nil {
		LogDebug("Error creating themes directory: %v", err)
		return nil, fmt.Errorf("error creating themes directory: %w", err)
	}

	// Read the directory
	entries, err := os.ReadDir(themesDir)
	if err != nil {
		LogDebug("Error reading themes directory: %v", err)
		return nil, fmt.Errorf("error reading themes directory: %w", err)
	}

	// Add each theme to the list
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			LogDebug("Found theme: %s", entry.Name())
			themes = append(themes, entry.Name())
		}
	}

	LogDebug("Found %d themes", len(themes))
	return themes, nil
}