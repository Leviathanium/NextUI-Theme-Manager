// src/internal/themes/common.go
// Common utilities for theme operations

package themes

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"nextui-themes/internal/logging"
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

	// Theme directories to create
	directories := []string{
		filepath.Join(cwd, "Themes"),
		filepath.Join(cwd, "Themes", "Imports"),
		filepath.Join(cwd, "Themes", "Exports"),
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
		filepath.Join(cwd, "Themes", "Imports", "README.txt"): `# Theme Import Directory

Place theme packages (directories with .theme extension) here to import them.
Themes should contain a manifest.json file and the appropriate theme files.`,

		filepath.Join(cwd, "Themes", "Exports", "README.txt"): `# Theme Export Directory

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