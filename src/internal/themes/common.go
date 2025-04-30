// internal/themes/common.go
package themes

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"thememanager/internal/logging"
	"thememanager/internal/system"
)

// Constants for max number of backups to keep
const MaxBackups = 3

// EnsureDirectoryStructure creates all necessary directories for the application
func EnsureDirectoryStructure() error {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Create required directories
	dirs := []string{
		filepath.Join(cwd, "Themes"),
		filepath.Join(cwd, "Overlays"),
		filepath.Join(cwd, "Backups", "Themes"),
		filepath.Join(cwd, "Backups", "Overlays"),
		filepath.Join(cwd, "Catalog", "Themes"),
		filepath.Join(cwd, "Catalog", "Overlays"),
		filepath.Join(cwd, "Logs"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			logging.LogDebug("Error creating directory %s: %v", dir, err)
			return fmt.Errorf("error creating directory %s: %w", dir, err)
		}
	}

	return nil
}

// IsThemeDownloaded checks if a theme is already downloaded
func IsThemeDownloaded(themeName string) bool {
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return false
	}

	themePath := filepath.Join(cwd, "Themes", themeName+system.ThemeExtension)
	_, err = os.Stat(themePath)
	return err == nil
}

// IsOverlayDownloaded checks if an overlay pack is already downloaded
func IsOverlayDownloaded(overlayName string) bool {
	cwd, err := os.Getwd()
	if err != nil {
		logging.LogDebug("Error getting current directory: %v", err)
		return false
	}

	overlayPath := filepath.Join(cwd, "Overlays", overlayName+system.OverlayExtension)
	_, err = os.Stat(overlayPath)
	return err == nil
}

// CopyDir recursively copies a directory
func CopyDir(src, dst string) error {
	logging.LogDebug("Copying directory: %s -> %s", src, dst)

	// Get file info
	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create destination directory with same permissions
	if err := os.MkdirAll(dst, info.Mode()); err != nil {
		return err
	}

	// Read source directory
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// Copy each entry
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursively copy directory
			if err := CopyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy file
			if err := CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// CopyFile copies a single file
func CopyFile(src, dst string) error {
	logging.LogDebug("Copying file: %s -> %s", src, dst)

	// Create destination directory if it doesn't exist
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dstDir, err)
	}

	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Get file info
	info, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get source file info: %w", err)
	}

	// Create destination file
	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	// Copy contents
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	logging.LogDebug("Successfully copied file: %s", dst)
	return nil
}

// GetThemePath returns the full path to a theme
func GetThemePath(themeName string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %w", err)
	}

	// Check if theme name already has extension
	if !strings.HasSuffix(themeName, system.ThemeExtension) {
		themeName = themeName + system.ThemeExtension
	}

	return filepath.Join(cwd, "Themes", themeName), nil
}

// GetOverlayPath returns the full path to an overlay
func GetOverlayPath(overlayName string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %w", err)
	}

	// Check if overlay name already has extension
	if !strings.HasSuffix(overlayName, system.OverlayExtension) {
		overlayName = overlayName + system.OverlayExtension
	}

	return filepath.Join(cwd, "Overlays", overlayName), nil
}

// GetCatalogThemePath returns the full path to a theme in the catalog
func GetCatalogThemePath(themeName string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %w", err)
	}

	// Check if theme name already has extension
	if !strings.HasSuffix(themeName, system.ThemeExtension) {
		themeName = themeName + system.ThemeExtension
	}

	return filepath.Join(cwd, "Catalog", "Themes", themeName), nil
}

// GetCatalogOverlayPath returns the full path to an overlay in the catalog
func GetCatalogOverlayPath(overlayName string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %w", err)
	}

	// Check if overlay name already has extension
	if !strings.HasSuffix(overlayName, system.OverlayExtension) {
		overlayName = overlayName + system.OverlayExtension
	}

	return filepath.Join(cwd, "Catalog", "Overlays", overlayName), nil
}

// GetThemeBackupPath returns the full path to a theme backup
func GetThemeBackupPath(backupName string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %w", err)
	}

	return filepath.Join(cwd, "Backups", "Themes", backupName), nil
}

// GetOverlayBackupPath returns the full path to an overlay backup
func GetOverlayBackupPath(backupName string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %w", err)
	}

	return filepath.Join(cwd, "Backups", "Overlays", backupName), nil
}

// EnsureMediaDirectories creates all necessary media directories
func EnsureMediaDirectories(systemPaths *system.SystemPaths) error {
	return system.EnsureMediaDirectories(systemPaths)
}

// CleanDirectory removes all files in a directory but keeps the directory itself
func CleanDirectory(dir string) error {
	logging.LogDebug("Cleaning directory: %s", dir)

	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		logging.LogDebug("Directory does not exist, creating it: %s", dir)
		return os.MkdirAll(dir, 0755)
	}

	// Read directory
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("error reading directory: %w", err)
	}

	// Remove each entry
	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())

		if entry.IsDir() {
			// Remove directory and all contents
			if err := os.RemoveAll(path); err != nil {
				logging.LogDebug("Error removing directory %s: %v", path, err)
				// Continue with other entries
			}
		} else {
			// Remove file
			if err := os.Remove(path); err != nil {
				logging.LogDebug("Error removing file %s: %v", path, err)
				// Continue with other entries
			}
		}
	}

	logging.LogDebug("Successfully cleaned directory: %s", dir)
	return nil
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// CreateEmptyFile creates an empty file
func CreateEmptyFile(path string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Create empty file
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	return nil
}