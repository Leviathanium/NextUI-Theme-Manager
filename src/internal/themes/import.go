// src/internal/themes/import.go
// Simplified implementation of theme import functionality

package themes

import (
	"fmt"
	"os"
	"path/filepath"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/system"
	"nextui-themes/internal/ui"
)

// ImportTheme imports a theme package
func ImportTheme(themeName string) error {
	// Create logging directory if it doesn't exist
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	logsDir := filepath.Join(cwd, "Logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("error creating logs directory: %w", err)
	}

	// Create logger
	logger := &Logger{
		DebugFn: logging.LogDebug,
	}

	logger.DebugFn("Starting theme import for: %s", themeName)

	// Full path to theme - look in Imports directory
	themePath := filepath.Join(cwd, "Themes", "Imports", themeName)

	// Validate theme
	manifest, err := ValidateTheme(themePath, logger)
	if err != nil {
		logger.DebugFn("Theme validation failed: %v", err)
		return fmt.Errorf("theme validation failed: %w", err)
	}

	// Import theme components based on path mappings
	if err := importThemeFiles(themePath, manifest, logger); err != nil {
		logger.DebugFn("Error importing theme files: %v", err)
		return fmt.Errorf("error importing theme files: %w", err)
	}

	logger.DebugFn("Theme import completed successfully: %s", themeName)

	// Show success message to user
	ui.ShowMessage(fmt.Sprintf("Theme '%s' by %s imported successfully!",
		manifest.ThemeInfo.Name, manifest.ThemeInfo.Author), "3")

	return nil
}

// importThemeFiles copies all files from the theme to the system based on path mappings
func importThemeFiles(themePath string, manifest *ThemeManifest, logger *Logger) error {
	// Get system paths
	systemPaths, err := system.GetSystemPaths()
	if err != nil {
		logger.DebugFn("Error getting system paths: %v", err)
		// Continue anyway with just the path mappings
	}

	// Ensure media directories exist
	if systemPaths != nil {
		if err := system.EnsureMediaDirectories(systemPaths); err != nil {
			logger.DebugFn("Warning: Failed to ensure media directories: %v", err)
		}
	}

	// Process wallpaper mappings
	for _, mapping := range manifest.PathMappings.Wallpapers {
		srcPath := filepath.Join(themePath, mapping.ThemePath)
		dstPath := mapping.SystemPath

		// Copy the file
		if err := copyMappedFile(srcPath, dstPath, logger); err != nil {
			logger.DebugFn("Warning: Failed to copy wallpaper: %v", err)
			// Continue with other files
		}
	}

	// Process icon mappings
	for _, mapping := range manifest.PathMappings.Icons {
		srcPath := filepath.Join(themePath, mapping.ThemePath)
		dstPath := mapping.SystemPath

		// Copy the file
		if err := copyMappedFile(srcPath, dstPath, logger); err != nil {
			logger.DebugFn("Warning: Failed to copy icon: %v", err)
			// Continue with other files
		}
	}

	// Process overlay mappings
	for _, mapping := range manifest.PathMappings.Overlays {
		srcPath := filepath.Join(themePath, mapping.ThemePath)
		dstPath := mapping.SystemPath

		// Copy the file
		if err := copyMappedFile(srcPath, dstPath, logger); err != nil {
			logger.DebugFn("Warning: Failed to copy overlay: %v", err)
			// Continue with other files
		}
	}

	// Process font mappings
	for fontType, mapping := range manifest.PathMappings.Fonts {
		srcPath := filepath.Join(themePath, mapping.ThemePath)
		dstPath := mapping.SystemPath

		// Copy the file
		if err := copyMappedFile(srcPath, dstPath, logger); err != nil {
			logger.DebugFn("Warning: Failed to copy font %s: %v", fontType, err)
			// Continue with other files
		}
	}

	// Process settings mappings
	for settingType, mapping := range manifest.PathMappings.Settings {
		srcPath := filepath.Join(themePath, mapping.ThemePath)
		dstPath := mapping.SystemPath

		// Copy the file
		if err := copyMappedFile(srcPath, dstPath, logger); err != nil {
			logger.DebugFn("Warning: Failed to copy setting %s: %v", settingType, err)
			// Continue with other files
		}
	}

	return nil
}

// copyMappedFile copies a file from source to destination with appropriate checks
func copyMappedFile(srcPath, dstPath string, logger *Logger) error {
	// Check if source file exists
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		logger.DebugFn("Source file does not exist: %s", srcPath)
		return fmt.Errorf("source file does not exist: %s", srcPath)
	}

	// Create destination directory
	dstDir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		logger.DebugFn("Failed to create destination directory: %v", err)
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Copy the file
	if err := CopyFile(srcPath, dstPath); err != nil {
		logger.DebugFn("Failed to copy file: %v", err)
		return fmt.Errorf("failed to copy file: %w", err)
	}

	logger.DebugFn("Copied file: %s -> %s", srcPath, dstPath)
	return nil
}