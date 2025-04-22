// src/internal/themes/common.go
// Common utilities for theme operations

package themes

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
    "regexp"
    "strings"
	"nextui-themes/internal/logging"
	"nextui-themes/internal/system"  // Add this import

)

// CopyFile copies a file from src to dst
func CopyFile(src, dst string) error {
	logging.LogDebug("Copying %s to %s", src, dst)

	// Create the destination directory if it doesn't exist
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		logging.LogDebug("Error creating directory %s: %v", dstDir, err)
		return fmt.Errorf("failed to create directory %s: %w", dstDir, err)
	}

	srcFile, err := os.Open(src)
	if err != nil {
		logging.LogDebug("Error opening source file: %v", err)
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		logging.LogDebug("Error creating destination file: %v", err)
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		logging.LogDebug("Error copying file: %v", err)
		return fmt.Errorf("failed to copy file: %w", err)
	}

	logging.LogDebug("Successfully copied %d bytes", bytes)
	return nil
}

// EnsureThemeDirectoryStructure creates all the necessary directories for theme management
func EnsureThemeDirectoryStructure() error {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return err
	}

	// Theme directories to create - now directly using Themes and Exports
	directories := []string{
		filepath.Join(cwd, "Themes"),
		filepath.Join(cwd, "Exports"),
		filepath.Join(cwd, "Logs"),
	}

	// Create each directory
	for _, dir := range directories {
		if err := os.MkdirAll(dir, 0755); err != nil {
			logging.LogDebug("Error creating directory %s: %v", dir, err)
			return err
		}
	}

	logging.LogDebug("Theme directory structure created")
	return nil
}

// CreatePlaceholderFiles creates README files in empty directories
func CreatePlaceholderFiles() error {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return err
	}

	// Define placeholder files
	placeholders := map[string]string{
		filepath.Join(cwd, "Themes", "README.txt"): `# Theme Directory

Place theme packages (directories with .theme extension) here to import them.
Themes should contain a manifest.json file and the appropriate theme files.`,

		filepath.Join(cwd, "Exports", "README.txt"): `# Theme Export Directory

Exported theme packages will be placed here with sequential names (theme_1.theme, theme_2.theme, etc.)`,
	}

	// Create each placeholder file if the directory is empty
	for filePath, content := range placeholders {
		dir := filepath.Dir(filePath)

		// Check if directory is empty (except for other README files)
		entries, err := os.ReadDir(dir)
		if err != nil {
			logging.LogDebug("Error reading directory %s: %v", dir, err)
			continue
		}

		hasContent := false
		for _, entry := range entries {
			if !entry.IsDir() && entry.Name() != "README.txt" {
				hasContent = true
				break
			}
		}

		// Create README if directory is empty or only contains README
		if !hasContent {
			if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
				logging.LogDebug("Error creating placeholder file %s: %v", filePath, err)
			}
		}
	}

	return nil
}

// GetSystemIconDestination determines the correct destination path for a system icon
// It ensures icons are named to match exact ROM directory names based on system tags
func GetSystemIconDestination(iconSrcPath, iconName, dstPath string, systemPaths *system.SystemPaths, logger *Logger) (string, error) {
	// Extract system tag from the filename
	tagRegex := regexp.MustCompile(`\((.*?)\)`)
	matches := tagRegex.FindStringSubmatch(iconName)

	if len(matches) < 2 || matches[1] == "" {
		// No tag found, use destination as-is
		logger.DebugFn("No system tag found in icon name: %s", iconName)
		return dstPath, nil
	}

	// Get system tag
	systemTag := matches[1]
	logger.DebugFn("Found system tag in icon: %s - Tag: %s", iconName, systemTag)

	// Handle special system icons
	if iconName == "Recently Played.png" ||
	   iconName == "Collections.png" ||
	   iconName == "Tools.png" {
		// These special icons don't need tag-based renaming
		logger.DebugFn("Special system icon, no renaming needed: %s", iconName)
		return dstPath, nil
	}

	// Look for exact ROM directory matching this tag
	var exactSystemName string
	var matchFound bool

	// First pass: look for exact matches
	for _, system := range systemPaths.Systems {
		if system.Tag == systemTag {
			exactSystemName = system.Name
			matchFound = true
			logger.DebugFn("Found exact ROM directory match for tag '%s': %s", systemTag, exactSystemName)
			break
		}
	}

	// If no match found, try case-insensitive matching
	if !matchFound {
		systemTagLower := strings.ToLower(systemTag)
		for _, system := range systemPaths.Systems {
			if strings.ToLower(system.Tag) == systemTagLower {
				exactSystemName = system.Name
				matchFound = true
				logger.DebugFn("Found case-insensitive ROM directory match for tag '%s': %s", systemTag, exactSystemName)
				break
			}
		}
	}

	// If we found a matching ROM directory, rename the icon to match it exactly
	if matchFound && exactSystemName != "" {
		// Get the media directory path
		mediaDir := filepath.Dir(dstPath)

		// Create the new destination path with the exact ROM directory name
		newDstPath := filepath.Join(mediaDir, exactSystemName + ".png")

		logger.DebugFn("Renaming system icon from '%s' to match ROM directory: '%s'", iconName, exactSystemName + ".png")
		return newDstPath, nil
	}

	// No matching ROM directory found, keep original destination
	logger.DebugFn("No matching ROM directory found for system tag '%s', using original name", systemTag)
	return dstPath, nil
}