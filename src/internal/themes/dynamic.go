// src/internal/themes/dynamic.go
// Dynamic theme operations

package themes

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"nextui-themes/internal/logging"
	"nextui-themes/internal/system"
)

// ThemeFile represents a background image in a dynamic theme
type ThemeFile struct {
	SourcePath string // Path to the source image
	TargetPath string // Path where it should be copied
}

// ListDynamicThemes returns a list of available dynamic themes
func ListDynamicThemes() ([]string, error) {
	var themes []string

	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return nil, fmt.Errorf("error getting current directory: %w", err)
	}

	// Get the dynamic themes directory (using absolute path) - updated to use Themes/Imports
	themesDir := filepath.Join(cwd, "Themes", "Imports")
	logging.LogDebug("Listing themes from directory: %s", themesDir)

	// Check if the directory exists
	_, err = os.Stat(themesDir)
	if os.IsNotExist(err) {
		logging.LogDebug("Dynamic themes directory does not exist: %s", themesDir)
		return nil, fmt.Errorf("dynamic themes directory does not exist: %s", themesDir)
	} else if err != nil {
		logging.LogDebug("Error checking dynamic themes directory: %v", err)
		return nil, fmt.Errorf("error checking themes directory: %w", err)
	}

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(themesDir, 0755); err != nil {
		logging.LogDebug("Error creating themes directory: %v", err)
		return nil, fmt.Errorf("error creating themes directory: %w", err)
	}

	// Read the directory
	entries, err := os.ReadDir(themesDir)
	if err != nil {
		logging.LogDebug("Error reading themes directory: %v", err)
		return nil, fmt.Errorf("error reading themes directory: %w", err)
	}

	// Add each theme to the list
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			logging.LogDebug("Found theme: %s", entry.Name())
			themes = append(themes, entry.Name())
		}
	}

	logging.LogDebug("Found %d themes", len(themes))
	return themes, nil
}

// ScanDynamicTheme scans a dynamic theme and returns files to copy
func ScanDynamicTheme(themeName string, systemPaths *system.SystemPaths) ([]ThemeFile, error) {
	var themeFiles []ThemeFile

	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting current directory: %w", err)
	}

	// Get the theme base directory (using absolute path) - updated to use Themes/Imports
	themeDir := filepath.Join(cwd, "Themes", "Imports", themeName)
	logging.LogDebug("Scanning dynamic theme directory: %s", themeDir)

	// Check if the theme directory exists
	_, err = os.Stat(themeDir)
	if os.IsNotExist(err) {
		logging.LogDebug("Theme directory does not exist: %s", themeDir)
		return nil, fmt.Errorf("theme directory does not exist: %s", themeDir)
	} else if err != nil {
		logging.LogDebug("Error checking theme directory: %v", err)
		return nil, fmt.Errorf("error checking theme directory: %w", err)
	}

	// Rest of the function remains the same...
	// Handle Root background
	rootBg := filepath.Join(themeDir, "Root", "bg.png")
	if _, err := os.Stat(rootBg); err == nil {
		logging.LogDebug("Found Root background: %s", rootBg)
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
		logging.LogDebug("Root background not found: %s, error: %v", rootBg, err)
	}

	// Rest of the scanning code remains unchanged...

	logging.LogDebug("Theme scan complete, found %d files", len(themeFiles))
	return themeFiles, nil
}

// ApplyDynamicTheme applies a dynamic theme
func ApplyDynamicTheme(themeName string) error {
	logging.LogDebug("Applying dynamic theme: %s", themeName)

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

	// Scan the dynamic theme
	themeFiles, err := ScanDynamicTheme(themeName, systemPaths)
	if err != nil {
		logging.LogDebug("Error scanning theme: %v", err)
		return fmt.Errorf("error scanning theme: %w", err)
	}

	logging.LogDebug("Found %d theme files to apply", len(themeFiles))

	// Apply each theme file
	for _, file := range themeFiles {
		logging.LogDebug("Copying %s to %s", file.SourcePath, file.TargetPath)

		// Create parent directories if needed
		parentDir := filepath.Dir(file.TargetPath)
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			logging.LogDebug("Error creating directory %s: %v", parentDir, err)
			return fmt.Errorf("error creating directory %s: %w", parentDir, err)
		}

		// Copy the file
		if err := CopyFile(file.SourcePath, file.TargetPath); err != nil {
			logging.LogDebug("Error copying %s to %s: %v", file.SourcePath, file.TargetPath, err)
			return fmt.Errorf("error copying %s to %s: %w", file.SourcePath, file.TargetPath, err)
		}
	}

	logging.LogDebug("Dynamic theme applied successfully")
	return nil
}
